package admin

import (
	"strings"
	"testing"
)

// TestIsValidEmail тестирует валидацию email адресов
func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		name  string
		email string
		want  bool
	}{
		// Валидные email
		{
			name:  "valid simple email",
			email: "user@example.com",
			want:  true,
		},
		{
			name:  "valid email with numbers",
			email: "user123@example.com",
			want:  true,
		},
		{
			name:  "valid email with dots",
			email: "first.last@example.com",
			want:  true,
		},
		{
			name:  "valid email with plus",
			email: "user+tag@example.com",
			want:  true,
		},
		{
			name:  "valid email with underscore",
			email: "user_name@example.com",
			want:  true,
		},
		{
			name:  "valid email with hyphen",
			email: "user-name@example-domain.com",
			want:  true,
		},
		{
			name:  "valid email with subdomain",
			email: "user@mail.example.com",
			want:  true,
		},

		// Невалидные email
		{
			name:  "empty string",
			email: "",
			want:  false,
		},
		{
			name:  "no @ symbol",
			email: "userexample.com",
			want:  false,
		},
		{
			name:  "no domain",
			email: "user@",
			want:  false,
		},
		{
			name:  "no local part",
			email: "@example.com",
			want:  false,
		},
		{
			name:  "no TLD",
			email: "user@example",
			want:  false,
		},
		{
			name:  "double @ symbol",
			email: "user@@example.com",
			want:  false,
		},
		{
			name:  "spaces in email",
			email: "user name@example.com",
			want:  false,
		},
		{
			name:  "invalid characters",
			email: "user#name@example.com",
			want:  false,
		},
		{
			name:  "too long (>254 chars)",
			email: strings.Repeat("a", 250) + "@example.com",
			want:  false,
		},
		{
			name:  "missing dot in domain",
			email: "user@examplecom",
			want:  false,
		},
		{
			name:  "TLD too short",
			email: "user@example.c",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isValidEmail(tt.email)
			if got != tt.want {
				t.Errorf("isValidEmail(%q) = %v, want %v", tt.email, got, tt.want)
			}
		})
	}
}

// TestValidatePassword тестирует валидацию паролей
func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
		errMsg   string
	}{
		// Валидные пароли
		{
			name:     "valid password with all requirements",
			password: "Password123",
			wantErr:  false,
		},
		{
			name:     "valid password longer",
			password: "MySecurePassword123",
			wantErr:  false,
		},
		{
			name:     "valid password at minimum length",
			password: "Pass123w",
			wantErr:  false,
		},
		{
			name:     "valid password at maximum length (72 chars)",
			password: "A1" + strings.Repeat("a", 69) + "B",
			wantErr:  false,
		},
		{
			name:     "valid password with special chars",
			password: "P@ssw0rd!#$",
			wantErr:  false,
		},

		// Невалидные пароли - слишком короткие
		{
			name:     "too short - 7 chars",
			password: "Pass123",
			wantErr:  true,
			errMsg:   "пароль должен содержать минимум 8 символов",
		},
		{
			name:     "too short - empty",
			password: "",
			wantErr:  true,
			errMsg:   "пароль должен содержать минимум 8 символов",
		},

		// Невалидные пароли - слишком длинные
		{
			name:     "too long - 73 chars",
			password: "A1" + strings.Repeat("a", 70) + "B",
			wantErr:  true,
			errMsg:   "пароль не может быть длиннее 72 символов",
		},
		{
			name:     "too long - 100 chars",
			password: strings.Repeat("Password123", 10),
			wantErr:  true,
			errMsg:   "пароль не может быть длиннее 72 символов",
		},

		// Невалидные пароли - нет заглавных букв
		{
			name:     "no uppercase",
			password: "password123",
			wantErr:  true,
			errMsg:   "пароль должен содержать заглавные и строчные буквы, а также цифры",
		},

		// Невалидные пароли - нет строчных букв
		{
			name:     "no lowercase",
			password: "PASSWORD123",
			wantErr:  true,
			errMsg:   "пароль должен содержать заглавные и строчные буквы, а также цифры",
		},

		// Невалидные пароли - нет цифр
		{
			name:     "no digits",
			password: "PasswordOnly",
			wantErr:  true,
			errMsg:   "пароль должен содержать заглавные и строчные буквы, а также цифры",
		},

		// Невалидные пароли - нет ничего
		{
			name:     "only numbers",
			password: "12345678",
			wantErr:  true,
			errMsg:   "пароль должен содержать заглавные и строчные буквы, а также цифры",
		},
		{
			name:     "only lowercase",
			password: "abcdefgh",
			wantErr:  true,
			errMsg:   "пароль должен содержать заглавные и строчные буквы, а также цифры",
		},
		{
			name:     "only uppercase",
			password: "ABCDEFGH",
			wantErr:  true,
			errMsg:   "пароль должен содержать заглавные и строчные буквы, а также цифры",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validatePassword(tt.password)

			if tt.wantErr {
				if err == nil {
					t.Errorf("validatePassword(%q) expected error but got nil", tt.password)
					return
				}
				if err.Error() != tt.errMsg {
					t.Errorf("validatePassword(%q) error = %q, want %q", tt.password, err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("validatePassword(%q) unexpected error = %v", tt.password, err)
				}
			}
		})
	}
}

// BenchmarkIsValidEmail бенчмарк для проверки производительности валидации email
func BenchmarkIsValidEmail(b *testing.B) {
	email := "user@example.com"
	for i := 0; i < b.N; i++ {
		isValidEmail(email)
	}
}

// BenchmarkValidatePassword бенчмарк для проверки производительности валидации пароля
func BenchmarkValidatePassword(b *testing.B) {
	password := "Password123"
	for i := 0; i < b.N; i++ {
		validatePassword(password)
	}
}
