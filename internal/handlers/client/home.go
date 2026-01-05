package client

import (
	"net/http"
	"prommsc/internal/handlers/shared"
	"prommsc/models"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	data := models.PageData{
		Title:   "ProminaYou - Главная",
		Content: "Добро пожаловать на наш сайт!",
	}
	shared.RenderTemplate(w, r, "index.html", data)
}
