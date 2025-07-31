package handlers

import (
	"net/http"
	"prommsc/models"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, "index.html", models.PageData{
		Title:   "ProminaYou - Главная",
		Content: "Добро пожаловать на наш сайт!",
	})
}
