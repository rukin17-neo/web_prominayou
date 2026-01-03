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
	services, err := h.repo.GetAll()
	if err != nil {
		http.Error(w, "Ошибка загрузки услуг", http.StatusInternalServerError)
		return
	}

	data := struct {
		Title    string
		Services []models.Service
	}{
		Title:    "Наши услуги",
		Services: services,
	}

	shared.RenderTemplate(w, r, "services.html", data)
}
