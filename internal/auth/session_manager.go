package auth

import (
	"log"
	"net/http"
	"prommsc/models"
	"strings"
	"time"
)

const (
	sessionCookieName  = "prommsc_session"
	rememberCookieName = "prommsc_remember"
	sessionDuration    = 24 * time.Hour       // 24 часа
	rememberDuration   = 30 * 24 * time.Hour  // 30 дней
)

type SessionManager struct {
	sessionRepo *models.SessionRepository
	userRepo    *models.UserRepository
}

func NewSessionManager(sessionRepo *models.SessionRepository, userRepo *models.UserRepository) *SessionManager {
	return &SessionManager{
		sessionRepo: sessionRepo,
		userRepo:    userRepo,
	}
}

// CreateSession создаёт новую сессию для пользователя
func (sm *SessionManager) CreateSession(w http.ResponseWriter, r *http.Request, userID int, rememberMe bool) error {
	// Генерация session ID
	sessionID, err := GenerateSessionID()
	if err != nil {
		return err
	}

	// Определение срока действия
	duration := sessionDuration
	if rememberMe {
		duration = rememberDuration
	}

	// Создание сессии в БД
	session := &models.Session{
		ID:        sessionID,
		UserID:    userID,
		ExpiresAt: time.Now().Add(duration),
		IPAddress: getIP(r),
		UserAgent: r.UserAgent(),
	}

	if err := sm.sessionRepo.Create(session); err != nil {
		return err
	}

	// Установка session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    sessionID,
		Path:     "/",
		MaxAge:   int(duration.Seconds()),
		HttpOnly: true,
		Secure:   false, // TODO: установить в true для production с HTTPS
		SameSite: http.SameSiteLaxMode,
	})

	// Если remember me - создать remember token
	if rememberMe {
		if err := sm.setRememberToken(w, userID); err != nil {
			log.Printf("Ошибка создания remember token: %v", err)
			// Не возвращаем ошибку, сессия уже создана
		}
	}

	return nil
}

// GetSessionUser получает пользователя по сессии из cookie
func (sm *SessionManager) GetSessionUser(r *http.Request) (*models.User, error) {
	// Попытка получить session cookie
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		return nil, err
	}

	// Получение сессии из БД
	session, err := sm.sessionRepo.GetByID(cookie.Value)
	if err != nil {
		return nil, err
	}

	// Получение пользователя
	user, err := sm.userRepo.GetByID(session.UserID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// DestroySession удаляет сессию пользователя
func (sm *SessionManager) DestroySession(w http.ResponseWriter, r *http.Request) error {
	// Получение session cookie
	cookie, err := r.Cookie(sessionCookieName)
	if err == nil {
		// Удаление сессии из БД
		sm.sessionRepo.Delete(cookie.Value)
	}

	// Очистка session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	// Очистка remember cookie
	http.SetCookie(w, &http.Cookie{
		Name:     rememberCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	return nil
}

// GetUserByRememberToken пытается восстановить сессию по remember token
func (sm *SessionManager) GetUserByRememberToken(r *http.Request) (*models.User, error) {
	cookie, err := r.Cookie(rememberCookieName)
	if err != nil {
		return nil, err
	}

	user, err := sm.userRepo.GetByRememberToken(cookie.Value)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// CleanupExpiredSessions удаляет истекшие сессии
func (sm *SessionManager) CleanupExpiredSessions() error {
	if err := sm.sessionRepo.DeleteExpired(); err != nil {
		log.Printf("Ошибка очистки сессий: %v", err)
		return err
	}
	return nil
}

// setRememberToken создаёт и устанавливает remember token
func (sm *SessionManager) setRememberToken(w http.ResponseWriter, userID int) error {
	token, err := GenerateSecureToken(32)
	if err != nil {
		return err
	}

	expiresAt := time.Now().Add(rememberDuration)
	if err := sm.userRepo.SetRememberToken(userID, token, expiresAt); err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     rememberCookieName,
		Value:    token,
		Path:     "/",
		MaxAge:   int(rememberDuration.Seconds()),
		HttpOnly: true,
		Secure:   false, // TODO: установить в true для production
		SameSite: http.SameSiteLaxMode,
	})

	return nil
}

// getIP извлекает IP адрес из запроса
func getIP(r *http.Request) string {
	// Проверка на прокси
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		parts := strings.Split(forwarded, ",")
		return strings.TrimSpace(parts[0])
	}

	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	return r.RemoteAddr
}
