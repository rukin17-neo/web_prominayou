package admin

import "regexp"

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
