package handlers

import (
	"log"
	"net/http"
	"prommsc/models"
)

func ReviewsHandler(w http.ResponseWriter, r *http.Request) {
	templatePath := "reviews.html"
	reviews, err := models.GetAllReviews()
	if err != nil {
		log.Printf("Ошибка получения отзывов: %v", err)
		http.Error(w, "Ошибка загрузки отзывов", http.StatusInternalServerError)
		return
	}

	data := struct {
		Title   string
		Reviews []models.Review
	}{
		Title:   "Отзывы клиентов",
		Reviews: reviews,
	}

	RenderTemplate(w, templatePath, data)
}
