package client

import (
	"net/http"
	"prommsc/internal/handlers/shared"
	"prommsc/models"
)

func ReviewsHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем замоканные отзывы из модели
	reviews, err := models.GetAllReviews()
	if err != nil {
		// В случае ошибки используем пустой список
		reviews = []models.Review{}
	}

	data := models.ReviewsData{
		PageData: models.PageData{
			Title:   "Отзывы клиентов",
			Content: "Что говорят о нас клиенты",
		},
		Reviews: reviews,
	}

	shared.RenderTemplate(w, r, "reviews.html", data)
}
