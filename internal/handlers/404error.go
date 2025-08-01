package handlers

import (
	"net/http"
	"prommsc/models"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	RenderTemplate(w, "404.html", models.PageData{
		Title: "Страница не найдена",
	})
}
