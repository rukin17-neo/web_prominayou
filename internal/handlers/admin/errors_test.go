package admin

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestLogAndRespondWithError проверяет корректность обработки ошибок
func TestLogAndRespondWithError(t *testing.T) {
	tests := []struct {
		name           string
		operation      string
		err            error
		userMessage    string
		statusCode     int
		wantStatusCode int
		wantBody       string
	}{
		{
			name:           "internal server error",
			operation:      "database query",
			err:            errors.New("connection timeout"),
			userMessage:    ErrMsgInternal,
			statusCode:     http.StatusInternalServerError,
			wantStatusCode: http.StatusInternalServerError,
			wantBody:       ErrMsgInternal,
		},
		{
			name:           "not found error",
			operation:      "fetch user",
			err:            errors.New("user not found in database"),
			userMessage:    ErrMsgNotFound,
			statusCode:     http.StatusNotFound,
			wantStatusCode: http.StatusNotFound,
			wantBody:       ErrMsgNotFound,
		},
		{
			name:           "unauthorized error",
			operation:      "check permissions",
			err:            errors.New("insufficient permissions"),
			userMessage:    ErrMsgUnauthorized,
			statusCode:     http.StatusForbidden,
			wantStatusCode: http.StatusForbidden,
			wantBody:       ErrMsgUnauthorized,
		},
		{
			name:           "bad request error",
			operation:      "parse form",
			err:            errors.New("invalid form data"),
			userMessage:    ErrMsgBadRequest,
			statusCode:     http.StatusBadRequest,
			wantStatusCode: http.StatusBadRequest,
			wantBody:       ErrMsgBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()

			logAndRespondWithError(rec, tt.operation, tt.err, tt.userMessage, tt.statusCode)

			// Проверяем статус код
			if rec.Code != tt.wantStatusCode {
				t.Errorf("Status code = %d, want %d", rec.Code, tt.wantStatusCode)
			}

			// Проверяем что тело ответа содержит пользовательское сообщение
			body := rec.Body.String()
			if !contains(body, tt.wantBody) {
				t.Errorf("Response body = %q, want to contain %q", body, tt.wantBody)
			}

			// Важно: проверяем что технические детали (из err.Error()) НЕ попали в ответ
			// Это защищает от утечки информации
			errDetails := tt.err.Error()
			if contains(body, errDetails) {
				t.Errorf("Response body contains technical error details %q, should only contain user message", errDetails)
			}
		})
	}
}

// TestErrorConstants проверяет наличие и непустоту констант сообщений об ошибках
func TestErrorConstants(t *testing.T) {
	constants := map[string]string{
		"ErrMsgInternal":     ErrMsgInternal,
		"ErrMsgNotFound":     ErrMsgNotFound,
		"ErrMsgUnauthorized": ErrMsgUnauthorized,
		"ErrMsgBadRequest":   ErrMsgBadRequest,
		"ErrMsgLoadFailed":   ErrMsgLoadFailed,
		"ErrMsgUpdateFailed": ErrMsgUpdateFailed,
		"ErrMsgCreateFailed": ErrMsgCreateFailed,
		"ErrMsgDeleteFailed": ErrMsgDeleteFailed,
	}

	for name, value := range constants {
		if value == "" {
			t.Errorf("Constant %s is empty", name)
		}
	}
}

// contains проверяет содержит ли строка подстроку
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
