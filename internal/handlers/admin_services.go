package handlers

import (
	"encoding/json"
	"net/http"
	"prommsc/models"
	"strconv"
)

type AdminServicesHandler struct {
	repo *models.ServiceRepository
}

func NewAdminServicesHandler(repo *models.ServiceRepository) *AdminServicesHandler {
	return &AdminServicesHandler{repo: repo}
}

func (h *AdminServicesHandler) UpdateService(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	var service models.Service
	if err := json.NewDecoder(r.Body).Decode(&service); err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	if err := h.repo.Update(&service); err != nil {
		http.Error(w, "Ошибка обновления", http.StatusInternalServerError)
		return
	}
	w.WriteHeader((http.StatusOK))
}

func (h *AdminServicesHandler) DeleteService(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}
	if err := h.repo.Delete(id); err != nil {
		http.Error(w, "Ошибка удаления", http.StatusInternalServerError)
		return
	}
	w.WriteHeader((http.StatusOK))
}

func (h *AdminServicesHandler) ListServices(w http.ResponseWriter, r *http.Request) {
	services, err := h.repo.GetAll()
	if err != nil {
		http.Error(w, "Ошибка загрузки услуг", http.StatusInternalServerError)
		return
	}

	if r.Header.Get("Accept") == "application/json" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(services)
		return
	}

	data := struct {
		Title    string
		Services []models.Service
	}{
		Title:    "Управление услугами",
		Services: services,
	}
	RenderTemplate(w, "admin/services.html", data)
}
