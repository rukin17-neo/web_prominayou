package client

import (
	"net/http"
	"prommsc/internal/handlers/shared"
	"prommsc/models"
)

func ContactsHandler(w http.ResponseWriter, r *http.Request) {
	data := models.ContactsData{
		PageData: models.PageData{
			Title:   "Контакты",
			Content: "Как нас найти",
		},
		Address: "г. Москва, Мантулинская, д. 9, к.1",
		Phone:   "+7 (966) 055-00-77",
	}

	shared.RenderTemplate(w, r, "contacts.html", data)
}
