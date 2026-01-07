package client

import (
	"log"
	"net/http"
	"prommsc/internal/handlers/shared"
	"prommsc/models"
)

type MastersPageData struct {
	models.PageData
	Masters    []models.Master
	Pagination models.PaginationResult
}

type MastersHandler struct {
	repo *models.MastersRepository
}

func NewMastersHandler(repo *models.MastersRepository) *MastersHandler {
	return &MastersHandler{repo: repo}
}

func (h *MastersHandler) List(w http.ResponseWriter, r *http.Request) {
	// Получение параметров пагинации из запроса
	paginationParams := models.NewPaginationParams(r)

	// Получение мастеров с пагинацией
	masters, pagination, err := h.repo.GetAllWithPagination(paginationParams)
	if err != nil {
		log.Printf("Ошибка загрузки мастеров: %v", err)
		http.Error(w, "Ошибка загрузки мастеров", http.StatusInternalServerError)
		return
	}

	data := MastersPageData{
		PageData:   models.PageData{Title: "Наши мастера", Content: "Наша команда профессионалов."},
		Masters:    masters,
		Pagination: pagination,
	}
	shared.RenderTemplate(w, r, "masters.html", data)
}

// фото из бд
func (h *MastersHandler) GetPhoto(w http.ResponseWriter, r *http.Request) {
	shared.ServePhoto(h.repo)(w, r)
}
