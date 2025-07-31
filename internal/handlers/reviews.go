package handlers

import (
	"net/http"
	"prommsc/models"
)

func ReviewsHandler(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, "reviews.html", models.PageData{
		Title:   "Отзывы",
		Content: "Отлично!",
	})
}
