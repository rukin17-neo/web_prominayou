package auth

import "golang.org/x/crypto/bcrypt"

const bcryptCost = 12

// HashPassword хэширует пароль используя bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	return string(bytes), err
}

// VerifyPassword проверяет соответствие пароля хэшу
func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
