package shared

import (
	"fmt"
	"net/http"
	"prommsc/models"
	"strconv"
	"strings"
)

// ServePhoto возвращает HTTP handler для отдачи фотографий мастеров из БД
func ServePhoto(repo *models.MastersRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// извлекаем id из url
		pathParts := strings.Split(r.URL.Path, "/")
		if len(pathParts) < 4 {
			http.Error(w, "Неверный URL", http.StatusBadRequest)
			return
		}

		idStr := pathParts[3]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Неверный ID", http.StatusBadRequest)
			return
		}

		master, err := repo.GetByID(id)
		if err != nil {
			http.Error(w, "Мастер не найден", http.StatusNotFound)
			return
		}

		if len(master.PhotoData) == 0 {
			http.Error(w, "Фотография не найдена", http.StatusNotFound)
			return
		}

		// устанавливаем заголовки для изображения
		contentType := "image/jpeg" // значение по умолчанию
		if master.PhotoType.Valid {
			contentType = master.PhotoType.String
		}
		w.Header().Set("Content-Type", contentType)
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(master.PhotoData)))
		w.Header().Set("Cache-Control", "public, max-age=31536000") // кэширование на 1 год

		// отправляем данные изображения
		w.Write(master.PhotoData)
	}
}
