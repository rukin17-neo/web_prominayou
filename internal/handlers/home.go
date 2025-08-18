package handlers

import (
	"net/http"
	"prommsc/models"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	data := models.PageData{
		Title:   "ProminaYou - Главная",
		Content: "Добро пожаловать на наш сайт!",
	}
	RenderTemplate(w, "index.html", data)
}
