package config

import (
	"fmt"
	"net/smtp"
	"os"
)

type EmailConfig struct {
	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
	FromEmail    string
	FromName     string
	AppURL       string
}

// GetEmailConfig загружает конфигурацию email из переменных окружения
func GetEmailConfig() *EmailConfig {
	return &EmailConfig{
		SMTPHost:     os.Getenv("SMTP_HOST"),
		SMTPPort:     os.Getenv("SMTP_PORT"),
		SMTPUsername: os.Getenv("SMTP_USERNAME"),
		SMTPPassword: os.Getenv("SMTP_PASSWORD"),
		FromEmail:    os.Getenv("SMTP_FROM_EMAIL"),
		FromName:     os.Getenv("SMTP_FROM_NAME"),
		AppURL:       os.Getenv("APP_URL"),
	}
}

// SendEmail отправляет email используя SMTP
func (c *EmailConfig) SendEmail(to, subject, body string) error {
	// Проверка конфигурации
	if c.SMTPHost == "" || c.SMTPPort == "" {
		return fmt.Errorf("SMTP не настроен: проверьте SMTP_HOST и SMTP_PORT в .env")
	}

	// Формирование сообщения
	from := c.FromEmail
	if from == "" {
		from = "noreply@prommsc.local"
	}

	message := []byte(
		"From: " + c.FromName + " <" + from + ">\r\n" +
			"To: " + to + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/plain; charset=UTF-8\r\n" +
			"\r\n" +
			body + "\r\n")

	// Настройка аутентификации
	auth := smtp.PlainAuth("", c.SMTPUsername, c.SMTPPassword, c.SMTPHost)

	// Отправка
	addr := c.SMTPHost + ":" + c.SMTPPort
	err := smtp.SendMail(addr, auth, from, []string{to}, message)
	if err != nil {
		return fmt.Errorf("ошибка отправки email: %w", err)
	}

	return nil
}

// SendPasswordResetEmail отправляет письмо со ссылкой для сброса пароля
func (c *EmailConfig) SendPasswordResetEmail(to, token string) error {
	appURL := c.AppURL
	if appURL == "" {
		appURL = "http://localhost:8002"
	}

	resetURL := fmt.Sprintf("%s/admin/reset-password?token=%s", appURL, token)
	subject := "Сброс пароля - ProminaYou Admin"
	body := fmt.Sprintf(`Здравствуйте!

Вы запросили сброс пароля для панели администратора ProminaYou.

Для сброса пароля перейдите по следующей ссылке:
%s

Ссылка действительна в течение 1 часа.

Если вы не запрашивали сброс пароля, просто проигнорируйте это письмо.

С уважением,
Команда ProminaYou`, resetURL)

	return c.SendEmail(to, subject, body)
}
