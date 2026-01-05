package client

import (
	"fmt"
	"net/http"
	"prommsc/internal/handlers/shared"
	"prommsc/models"
	"strconv"
	"strings"
)

type MastersPageData struct {
	models.PageData
	Masters []models.Master
}

type MastersHandler struct {
	repo *models.MastersRepository
}

func NewMastersHandler(repo *models.MastersRepository) *MastersHandler {
	return &MastersHandler{repo: repo}
}

func (h *MastersHandler) List(w http.ResponseWriter, r *http.Request) {
	masters, err := h.repo.GetAll()
	if err != nil {
		http.Error(w, "Ошибка загрузки мастеров", http.StatusInternalServerError)
		return
	}

	data := MastersPageData{
		PageData: models.PageData{Title: "Наши мастера", Content: "Наша команда профессионалов."},
		Masters:  masters,
	}
	shared.RenderTemplate(w, r, "masters.html", data)
}

// фото из бд
func (h *MastersHandler) GetPhoto(w http.ResponseWriter, r *http.Request) {
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

	master, err := h.repo.GetByID(id)
	if err != nil {
		http.Error(w, "Мастер не найден", http.StatusNotFound)
		return
	}

	if len(master.PhotoData) == 0 {
		http.Error(w, "Фотография не найдена", http.StatusNotFound)
		return
	}

	// устанавливаем заголовки для изображения
	w.Header().Set("Content-Type", master.PhotoType)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(master.PhotoData)))
	w.Header().Set("Cache-Control", "public, max-age=31536000") // кэширование 1 год

	w.Write(master.PhotoData)
}
