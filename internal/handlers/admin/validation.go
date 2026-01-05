package admin

import (
	"errors"
	"regexp"
	"unicode"
)

// emailRegex - регулярное выражение для валидации email адресов
// RFC 5322 compliant (упрощенная версия для практического использования)
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// isValidEmail проверяет, является ли строка валидным email адресом
// Использует regex валидацию вместо простой проверки на наличие "@"
func isValidEmail(email string) bool {
	if email == "" {
		return false
	}

	// Проверка длины (RFC 5321: максимум 254 символа для email)
	if len(email) > 254 {
		return false
	}

	// Проверка формата через регулярное выражение
	return emailRegex.MatchString(email)
}

// validatePassword проверяет пароль на соответствие требованиям безопасности
// Возвращает nil если пароль валиден, иначе error с описанием проблемы
func validatePassword(password string) error {
	// Проверка минимальной длины
	if len(password) < 8 {
		return errors.New("пароль должен содержать минимум 8 символов")
	}

	// Проверка максимальной длины (bcrypt принимает максимум 72 байта)
	if len(password) > 72 {
		return errors.New("пароль не может быть длиннее 72 символов")
	}

	var (
		hasUpper  bool
		hasLower  bool
		hasDigit  bool
	)

	// Проверка сложности пароля
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		}
	}

	// Требуем наличие заглавных, строчных букв и цифр
	if !hasUpper || !hasLower || !hasDigit {
		return errors.New("пароль должен содержать заглавные и строчные буквы, а также цифры")
	}

	return nil
}
