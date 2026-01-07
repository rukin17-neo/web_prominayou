package admin

import (
	"encoding/json"
	"log"
	"net/http"
	"prommsc/internal/handlers"
	"prommsc/internal/handlers/shared"
	"prommsc/models"
	"strconv"

	"github.com/gorilla/mux"
)

type AdminServicesHandler struct {
	repo *models.ServiceRepository
}

func NewAdminServicesHandler(repo *models.ServiceRepository) *AdminServicesHandler {
	return &AdminServicesHandler{repo: repo}
}

// форма создания услуг
func (h *AdminServicesHandler) CreateServiceForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	shared.RenderTemplate(w, r, "admin/create_service.html", struct {
		Title       string
		CurrentUser *models.User
	}{
		Title:       "Добавление новой услуги",
		CurrentUser: handlers.GetCurrentUser(r),
	})
}

// форма редактирования
func (h *AdminServicesHandler) UpdateServiceForm(w http.ResponseWriter, r *http.Request) {
	log.Println("--- UpdateServiceForm вызван ---")
	vars := mux.Vars(r)
	log.Printf("Параметры URL: %+v", vars)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Ошибка преобразования ID: %v", err)
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}
	log.Printf("Получен ID: %d", id)

	service, err := h.repo.GetById(id)
	if err != nil {
		http.Error(w, "Услуга не найдена", http.StatusNotFound)
		return
	}

	shared.RenderTemplate(w, r, "admin/edit_service.html", struct {
		Title       string
		Service     *models.Service
		CurrentUser *models.User
	}{
		Title:       "Редактирование услуги",
		Service:     service,
		CurrentUser: handlers.GetCurrentUser(r),
	})
}

func (h *AdminServicesHandler) ListServices(w http.ResponseWriter, r *http.Request) {
	// Получение параметров пагинации из запроса
	paginationParams := models.NewPaginationParams(r)

	// Получение услуг с пагинацией
	services, pagination, err := h.repo.GetAllWithPagination(paginationParams)
	if err != nil {
		logAndRespondWithError(w, "GetAllServices", err, ErrMsgLoadFailed, http.StatusInternalServerError)
		return
	}

	if r.Header.Get("Accept") == "application/json" {
		w.Header().Set("Content-Type", "application/json")
		response := struct {
			Services   []models.Service         `json:"services"`
			Pagination models.PaginationResult `json:"pagination"`
		}{
			Services:   services,
			Pagination: pagination,
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	var editService *models.Service
	if editParam := r.URL.Query().Get("edit"); editParam != "" {
		if id, err := strconv.Atoi(editParam); err == nil {
			if s, getErr := h.repo.GetById(id); getErr == nil {
				editService = s
			}
		}
	}

	shared.RenderTemplate(w, r, "admin/services.html", struct {
		Title       string
		Services    []models.Service
		EditService *models.Service
		CurrentUser *models.User
		Pagination  models.PaginationResult
	}{
		Title:       "Управление услугами",
		Services:    services,
		EditService: editService,
		CurrentUser: handlers.GetCurrentUser(r),
		Pagination:  pagination,
	})
}

// form-data и json
func (h *AdminServicesHandler) CreateService(w http.ResponseWriter, r *http.Request) {
	log.Println("CreateService called")
	log.Printf("Method: %s", r.Method)
	log.Printf("Content-Type: %s", r.Header.Get("Content-Type"))

	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Ошибка парсинга формы", http.StatusBadRequest)
		return
	}

	// если есть id - редактируем
	if idStr := r.FormValue("id"); idStr != "" {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Неверный ID", http.StatusBadRequest)
			return
		}

		price, err := strconv.ParseFloat(r.FormValue("price"), 64)
		if err != nil {
			http.Error(w, "Неверный формат цены", http.StatusBadRequest)
			return
		}

		service := models.Service{
			ID:       id,
			Name:     r.FormValue("name"),
			Price:    price,
			Duration: r.FormValue("duration"),
		}

		if err := h.repo.Update(&service); err != nil {
			logAndRespondWithError(w, "UpdateService", err, ErrMsgUpdateFailed, http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/admin/services", http.StatusSeeOther)
		return
	}

	var service models.Service
	if r.Header.Get("Content-Type") == "application/json" {
		if err := json.NewDecoder(r.Body).Decode(&service); err != nil {
			http.Error(w, "Неверный формат данных", http.StatusBadRequest)
			return
		}
	} else {
		price, err := strconv.ParseFloat(r.FormValue("price"), 64)
		if err != nil {
			http.Error(w, "Неверный формат цены", http.StatusBadRequest)
			return
		}

		service = models.Service{
			Name:     r.FormValue("name"),
			Price:    price,
			Duration: r.FormValue("duration"),
		}
	}

	if err := h.repo.Create(&service); err != nil {
		logAndRespondWithError(w, "CreateService", err, ErrMsgCreateFailed, http.StatusInternalServerError)
		return
	}

	if r.Header.Get("Accept") == "application/json" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(service)
	} else {
		http.Redirect(w, r, "/admin/services", http.StatusSeeOther)
	}
}

func (h *AdminServicesHandler) UpdateService(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	if r.Method == http.MethodGet {
		service, err := h.repo.GetById(id)
		if err != nil {
			http.Error(w, "Услуга не найдена", http.StatusNotFound)
			return
		}

		data := struct {
			Title   string
			Service models.Service
		}{
			Title:   "Редактирование услуги",
			Service: *service,
		}
		shared.RenderTemplate(w, r, "admin/edit_service.html", data)
		return
	}

	var service models.Service
	if r.Header.Get("Content-Type") == "application/json" {
		if err := json.NewDecoder(r.Body).Decode(&service); err != nil {
			http.Error(w, "Неверный формат данных", http.StatusBadRequest)
			return
		}
	} else {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Ошибка парсинга формы", http.StatusBadRequest)
			return
		}

		price, err := strconv.ParseFloat(r.FormValue("price"), 64)
		if err != nil {
			http.Error(w, "Неверный формат цены", http.StatusBadRequest)
			return
		}

		service = models.Service{
			ID:       id,
			Name:     r.FormValue("name"),
			Price:    price,
			Duration: r.FormValue("duration"),
		}
	}

	if err := h.repo.Update(&service); err != nil {
		http.Error(w, "Ошибка обновления: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Header.Get("Accept") == "application/json" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(service)
	} else {
		http.Redirect(w, r, "/admin/services", http.StatusSeeOther)
	}
}

func (h *AdminServicesHandler) DeleteService(w http.ResponseWriter, r *http.Request) {
	// Получаем ID из URL параметров
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	if err := h.repo.Delete(id); err != nil {
		logAndRespondWithError(w, "DeleteService", err, ErrMsgDeleteFailed, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/services", http.StatusSeeOther)
}
