package client

import (
	"net/http"
	"prommsc/internal/handlers/shared"
	"prommsc/models"
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
	shared.ServePhoto(h.repo)(w, r)
}
