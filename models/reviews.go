package models

import "time"

type Review struct {
	ID      int
	Author  string
	Text    string
	Rating  int
	Created time.Time
}

// замоканные отзывы
func GetAllReviews() ([]Review, error) {
	reviews := []Review{
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

	return reviews, nil
}
