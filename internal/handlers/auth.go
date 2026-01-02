package handlers

import (
	"context"
	"log"
	"net/http"
	"prommsc/internal/auth"
	"prommsc/models"
)

type contextKey string

const UserContextKey contextKey = "user"

var sessionManager *auth.SessionManager

// SetSessionManager устанавливает session manager для middleware
func SetSessionManager(sm *auth.SessionManager) {
	sessionManager = sm
}

// AuthMiddleware проверяет аутентификацию пользователя
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Попытка получить пользователя по session cookie
		user, err := sessionManager.GetSessionUser(r)

		// 2. Если сессия не найдена, попробовать remember token
		if err != nil {
			user, err = sessionManager.GetUserByRememberToken(r)
			if err == nil && user != nil {
				// Создаём новую сессию для пользователя с remember token
				if err := sessionManager.CreateSession(w, r, user.ID, false); err != nil {
					log.Printf("Ошибка создания сессии из remember token: %v", err)
				}
			}
		}

		// 3. Если пользователь не найден - редирект на логин
		if user == nil {
			http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
			return
		}

		// 4. Добавляем пользователя в контекст запроса
		ctx := context.WithValue(r.Context(), UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetCurrentUser получает текущего пользователя из контекста запроса
func GetCurrentUser(r *http.Request) *models.User {
	if user, ok := r.Context().Value(UserContextKey).(*models.User); ok {
		return user
	}
	return nil
}
