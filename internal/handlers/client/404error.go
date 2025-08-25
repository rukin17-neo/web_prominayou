package client

import (
	"net/http"
	"prommsc/internal/handlers/shared"
	"prommsc/models"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	shared.RenderTemplate(w, "404.html", models.PageData{
		Title: "Страница не найдена",
	})
}
