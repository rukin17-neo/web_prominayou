package admin

import (
	"log"
	"net/http"
)

// logAndRespondWithError логирует детальную ошибку на сервере
// и отправляет клиенту общее сообщение без технических деталей
func logAndRespondWithError(w http.ResponseWriter, operation string, err error, userMessage string, statusCode int) {
	// Логируем детальную ошибку на сервере (включая err.Error())
	log.Printf("Ошибка [%s]: %v", operation, err)

	// Отправляем клиенту общее сообщение без технических деталей
	http.Error(w, userMessage, statusCode)
}

// Константы для стандартных сообщений об ошибках
const (
	ErrMsgInternal      = "Внутренняя ошибка сервера. Пожалуйста, попробуйте позже."
	ErrMsgNotFound      = "Запрашиваемый ресурс не найден."
	ErrMsgUnauthorized  = "Доступ запрещен."
	ErrMsgBadRequest    = "Неверный формат запроса."
	ErrMsgLoadFailed    = "Не удалось загрузить данные. Пожалуйста, попробуйте позже."
	ErrMsgUpdateFailed  = "Не удалось обновить данные. Пожалуйста, попробуйте позже."
	ErrMsgCreateFailed  = "Не удалось создать запись. Пожалуйста, попробуйте позже."
	ErrMsgDeleteFailed  = "Не удалось удалить запись. Пожалуйста, попробуйте позже."
)
