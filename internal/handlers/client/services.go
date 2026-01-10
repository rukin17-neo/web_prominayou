package client

import (
	"net/http"
	"prommsc/internal/handlers/shared"
	"prommsc/models"
)

type ServicesHandler struct {
	repo *models.ServiceRepository
}

func NewServicesHandler(repo *models.ServiceRepository) *ServicesHandler {
	return &ServicesHandler{repo: repo}
}

func (h *ServicesHandler) GetAllServices(w http.ResponseWriter, r *http.Request) {
	// Получение параметров пагинации из запроса
	paginationParams := models.NewPaginationParams(r)
	println("DEBUG: page=", paginationParams.Page, "limit=", paginationParams.Limit)

	// Получение услуг с пагинацией
	services, pagination, err := h.repo.GetAllWithPagination(paginationParams)
	if err != nil {
		println("DEBUG ERROR:", err.Error())
		http.Error(w, "Ошибка загрузки услуг", http.StatusInternalServerError)
		return
	}
	println("DEBUG: Got", len(services), "services, total=", pagination.Total)

	data := struct {
		Title      string
		Services   []models.Service
		Pagination models.PaginationResult
	}{
		Title:      "Наши услуги",
		Services:   services,
		Pagination: pagination,
	}

	shared.RenderTemplate(w, r, "services.html", data)
}
