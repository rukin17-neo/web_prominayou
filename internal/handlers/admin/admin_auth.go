package admin

import (
	"log"
	"net/http"
	"prommsc/config"
	"prommsc/internal/auth"
	"prommsc/internal/handlers/shared"
	"prommsc/models"
	"strings"
	"time"
)

type AdminAuthHandler struct {
	userRepo       *models.UserRepository
	sessionManager *auth.SessionManager
}

func NewAdminAuthHandler(userRepo *models.UserRepository, sessionManager *auth.SessionManager) *AdminAuthHandler {
	return &AdminAuthHandler{
		userRepo:       userRepo,
		sessionManager: sessionManager,
	}
}

// LoginPage отображает страницу входа
func (h *AdminAuthHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	type pageData struct {
		Title   string
		Error   string
		Success string
	}

	data := pageData{
		Title:   "Вход",
		Success: r.URL.Query().Get("success"),
	}

	shared.RenderTemplate(w, r, "admin/login.html", data)
}

// Login обрабатывает вход пользователя
func (h *AdminAuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Ошибка парсинга формы", http.StatusBadRequest)
		return
	}

	username := strings.TrimSpace(r.FormValue("username"))
	password := r.FormValue("password")
	rememberMe := r.FormValue("remember_me") == "1"

	// Валидация
	if username == "" || password == "" {
		h.renderLoginError(w, r, "Все поля обязательны")
		return
	}

	// Получение пользователя
	user, err := h.userRepo.GetByUsername(username)
	if err != nil {
		log.Printf("Неудачная попытка входа: пользователь=%s, IP=%s", username, r.RemoteAddr)
		h.renderLoginError(w, r, "Неверные учетные данные")
		return
	}

	// Проверка пароля
	if err := auth.VerifyPassword(user.PasswordHash, password); err != nil {
		log.Printf("Неудачная попытка входа: пользователь=%s, IP=%s", username, r.RemoteAddr)
		h.renderLoginError(w, r, "Неверные учетные данные")
		return
	}

	// Создание сессии
	if err := h.sessionManager.CreateSession(w, r, user.ID, rememberMe); err != nil {
		log.Printf("Ошибка создания сессии: %v", err)
		http.Error(w, "Ошибка создания сессии", http.StatusInternalServerError)
		return
	}

	log.Printf("Успешный вход: пользователь=%s, IP=%s", username, r.RemoteAddr)
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

// Logout обрабатывает выход пользователя
func (h *AdminAuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if err := h.sessionManager.DestroySession(w, r); err != nil {
		log.Printf("Ошибка при выходе: %v", err)
	}
	http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
}

// ForgotPasswordPage отображает страницу восстановления пароля
func (h *AdminAuthHandler) ForgotPasswordPage(w http.ResponseWriter, r *http.Request) {
	type pageData struct {
		Title   string
		Success bool
	}

	data := pageData{
		Title: "Восстановление пароля",
	}

	shared.RenderTemplate(w, r, "admin/forgot_password.html", data)
}

// ForgotPassword обрабатывает запрос на восстановление пароля
func (h *AdminAuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Ошибка парсинга формы", http.StatusBadRequest)
		return
	}

	email := strings.TrimSpace(r.FormValue("email"))

	// Валидация email (базовая)
	if email == "" || !strings.Contains(email, "@") {
		h.renderForgotPasswordSuccess(w, r) // Показываем success даже при ошибке - security
		return
	}

	// Получение пользователя
	user, err := h.userRepo.GetByEmail(email)
	if err != nil {
		// Не сообщаем, что email не найден - security best practice
		log.Printf("Запрос сброса пароля для несуществующего email: %s", email)
		h.renderForgotPasswordSuccess(w, r)
		return
	}

	// Генерация reset token
	token, err := auth.GenerateSecureToken(32)
	if err != nil {
		log.Printf("Ошибка генерации токена: %v", err)
		http.Error(w, "Ошибка генерации токена", http.StatusInternalServerError)
		return
	}

	// Сохранение токена (действует 1 час)
	expiresAt := time.Now().Add(1 * time.Hour)
	if err := h.userRepo.SetResetToken(user.ID, token, expiresAt); err != nil {
		log.Printf("Ошибка сохранения токена: %v", err)
		http.Error(w, "Ошибка сохранения токена", http.StatusInternalServerError)
		return
	}

	// Отправка email
	emailConfig := config.GetEmailConfig()
	if err := emailConfig.SendPasswordResetEmail(email, token); err != nil {
		log.Printf("Ошибка отправки email: %v", err)
		// Не показываем ошибку пользователю
	}

	log.Printf("Запрос сброса пароля: email=%s, IP=%s", email, r.RemoteAddr)
	h.renderForgotPasswordSuccess(w, r)
}

// ResetPasswordPage отображает страницу сброса пароля
func (h *AdminAuthHandler) ResetPasswordPage(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")

	type pageData struct {
		Title string
		Token string
		Error string
	}

	// Проверка токена
	if token == "" {
		data := pageData{
			Title: "Сброс пароля",
			Error: "Токен не указан",
		}
		shared.RenderTemplate(w, r, "admin/reset_password.html", data)
		return
	}

	// Валидация токена в БД
	_, err := h.userRepo.GetByResetToken(token)
	if err != nil {
		data := pageData{
			Title: "Сброс пароля",
			Error: "Неверный или истекший токен",
		}
		shared.RenderTemplate(w, r, "admin/reset_password.html", data)
		return
	}

	// Токен валиден
	data := pageData{
		Title: "Сброс пароля",
		Token: token,
	}
	shared.RenderTemplate(w, r, "admin/reset_password.html", data)
}

// ResetPassword обрабатывает сброс пароля
func (h *AdminAuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Ошибка парсинга формы", http.StatusBadRequest)
		return
	}

	token := r.FormValue("token")
	password := r.FormValue("password")
	passwordConfirm := r.FormValue("password_confirm")

	// Валидация
	if token == "" || password == "" || passwordConfirm == "" {
		h.renderResetPasswordError(w, r, token, "Все поля обязательны")
		return
	}

	if password != passwordConfirm {
		h.renderResetPasswordError(w, r, token, "Пароли не совпадают")
		return
	}

	if len(password) < 8 {
		h.renderResetPasswordError(w, r, token, "Пароль должен содержать минимум 8 символов")
		return
	}

	// Получение пользователя по токену
	user, err := h.userRepo.GetByResetToken(token)
	if err != nil {
		h.renderResetPasswordError(w, r, token, "Неверный или истекший токен")
		return
	}

	// Хэширование нового пароля
	passwordHash, err := auth.HashPassword(password)
	if err != nil {
		log.Printf("Ошибка хэширования пароля: %v", err)
		http.Error(w, "Ошибка хэширования пароля", http.StatusInternalServerError)
		return
	}

	// Обновление пароля
	if err := h.userRepo.UpdatePassword(user.ID, passwordHash); err != nil {
		log.Printf("Ошибка обновления пароля: %v", err)
		http.Error(w, "Ошибка обновления пароля", http.StatusInternalServerError)
		return
	}

	// Очистка reset token
	if err := h.userRepo.ClearResetToken(user.ID); err != nil {
		log.Printf("Ошибка очистки токена: %v", err)
	}

	log.Printf("Пароль изменен: пользователь=%s", user.Username)

	// Редирект на страницу входа с сообщением об успехе
	http.Redirect(w, r, "/admin/login?success="+
		"Пароль успешно изменен. Войдите с новым паролем.", http.StatusSeeOther)
}

// Вспомогательные функции для рендера ошибок

func (h *AdminAuthHandler) renderLoginError(w http.ResponseWriter, r *http.Request, errorMsg string) {
	type pageData struct {
		Title string
		Error string
	}
	data := pageData{
		Title: "Вход",
		Error: errorMsg,
	}
	shared.RenderTemplate(w, r, "admin/login.html", data)
}

func (h *AdminAuthHandler) renderForgotPasswordSuccess(w http.ResponseWriter, r *http.Request) {
	type pageData struct {
		Title   string
		Success bool
	}
	data := pageData{
		Title:   "Восстановление пароля",
		Success: true,
	}
	shared.RenderTemplate(w, r, "admin/forgot_password.html", data)
}

func (h *AdminAuthHandler) renderResetPasswordError(w http.ResponseWriter, r *http.Request, token, errorMsg string) {
	type pageData struct {
		Title string
		Token string
		Error string
	}
	data := pageData{
		Title: "Сброс пароля",
		Token: token,
		Error: errorMsg,
	}
	shared.RenderTemplate(w, r, "admin/reset_password.html", data)
}
