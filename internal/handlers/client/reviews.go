package client

import (
	"net/http"
	"prommsc/internal/handlers/shared"
	"prommsc/models"
)

func ReviewsHandler(w http.ResponseWriter, r *http.Request) {
	// Получение параметров пагинации из запроса
	paginationParams := models.NewPaginationParams(r)

	// Получаем замоканные отзывы с пагинацией
	reviews, pagination, err := models.GetAllReviewsWithPagination(paginationParams)
	if err != nil {
		// В случае ошибки используем пустой список
		reviews = []models.Review{}
		pagination = models.PaginationResult{}
	}

	type ReviewsPageData struct {
		models.PageData
		Reviews    []models.Review
		Pagination models.PaginationResult
	}

	data := ReviewsPageData{
		PageData: models.PageData{
			Title:   "Отзывы клиентов",
			Content: "Что говорят о нас клиенты",
		},
		Reviews:    reviews,
		Pagination: pagination,
	}

	shared.RenderTemplate(w, r, "reviews.html", data)
}
