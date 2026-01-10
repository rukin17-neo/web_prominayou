package models

import "time"

type Review struct {
	ID      int
	Author  string
	Text    string
	Rating  int
	Created time.Time
}

// getAllReviewsMock возвращает все замоканные отзывы
func getAllReviewsMock() []Review {
	return []Review{
		{
			ID:      1,
			Author:  "Иван Петров",
			Text:    "Отличный сервис!",
			Rating:  5,
			Created: time.Now().Add(-24 * time.Hour), // 1 день назад
		},
		{
			ID:      2,
			Author:  "Мария Сидорова",
			Text:    "Очень довольна результатом, всем рекомендую!",
			Rating:  5,
			Created: time.Now().Add(-72 * time.Hour), // 3 дня назад
		},
		{
			ID:      3,
			Author:  "Алексей Иванов",
			Text:    "Хорошая работа.",
			Rating:  4,
			Created: time.Now().Add(-168 * time.Hour), // 1 неделя назад
		},
	}
}

// GetAllReviews возвращает все отзывы без пагинации (для обратной совместимости)
func GetAllReviews() ([]Review, error) {
	return getAllReviewsMock(), nil
}

// GetAllReviewsWithPagination возвращает отзывы с пагинацией
func GetAllReviewsWithPagination(params PaginationParams) ([]Review, PaginationResult, error) {
	allReviews := getAllReviewsMock()
	total := len(allReviews)

	// Вычисление индексов для среза
	start := params.GetOffset()
	end := start + params.Limit

	// Проверка границ
	if start >= total {
		return []Review{}, NewPaginationResult(total, params), nil
	}
	if end > total {
		end = total
	}

	// Получение среза отзывов
	reviews := allReviews[start:end]
	pagination := NewPaginationResult(total, params)

	return reviews, pagination, nil
}
