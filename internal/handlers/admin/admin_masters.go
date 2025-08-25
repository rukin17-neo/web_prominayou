package admin

import (
	"fmt"
	"io"
	"net/http"
	"prommsc/internal/handlers/shared"
	"prommsc/models"
	"strconv"
	"strings"
)

type AdminMastersHandler struct {
	repo *models.MastersRepository
}

func NewAdminMastersHandler(repo *models.MastersRepository) *AdminMastersHandler {
	return &AdminMastersHandler{repo: repo}
}

func (h *AdminMastersHandler) List(w http.ResponseWriter, r *http.Request) {
	masters, err := h.repo.GetAll()
	if err != nil {
		http.Error(w, "Ошибка загрузки мастеров: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var editMaster *models.Master
	// если есть параметр ?edit={id} - получаем мастера для редактирования
	if editID := r.URL.Query().Get("edit"); editID != "" {
		if id, err := strconv.Atoi(editID); err == nil {
			if m, e := h.repo.GetByID(id); e == nil {
				editMaster = m
			}
		}
	}

	type pageData struct {
		Title      string
		Masters    []models.Master
		EditMaster *models.Master
	}
	shared.RenderTemplate(w, "admin/masters.html", pageData{
		Title:      "Мастера",
		Masters:    masters,
		EditMaster: editMaster,
	})
}

func (h *AdminMastersHandler) CreateOrUpdate(w http.ResponseWriter, r *http.Request) {
	// парсинг формы с лимитом 16 мб
	if err := r.ParseMultipartForm(16 << 20); err != nil {
		http.Error(w, "Ошибка парсинга формы", http.StatusBadRequest)
		return
	}

	firstName := strings.TrimSpace(r.FormValue("first_name"))
	lastName := strings.TrimSpace(r.FormValue("last_name"))
	idStr := r.FormValue("id")

	// получение файла
	var photoData []byte
	var photoType string
	var photoURL string

	file, header, err := r.FormFile("photo")
	if err == nil && file != nil {
		defer file.Close()

		// читаем данные файла в память
		photoData, err = io.ReadAll(file)
		if err != nil {
			http.Error(w, "Ошибка чтения файла", http.StatusInternalServerError)
			return
		}

		// определяем MIME тип
		photoType = header.Header.Get("Content-Type")
		if photoType == "" {
			// определяем тип по расширению
			ext := strings.ToLower(header.Filename)
			switch {
			case strings.HasSuffix(ext, ".jpg"), strings.HasSuffix(ext, ".jpeg"):
				photoType = "image/jpeg"
			case strings.HasSuffix(ext, ".png"):
				photoType = "image/png"
			case strings.HasSuffix(ext, ".webp"):
				photoType = "image/webp"
			case strings.HasSuffix(ext, ".gif"):
				photoType = "image/gif"
			default:
				photoType = "image/jpeg"
			}
		}

		// валидация размера файла (максимум 5MB)
		if len(photoData) > 5*1024*1024 {
			http.Error(w, "Файл слишком большой. Максимальный размер: 5MB", http.StatusBadRequest)
			return
		}

		// валидация типа файла
		if !strings.HasPrefix(photoType, "image/") {
			http.Error(w, "Поддерживаются только изображения", http.StatusBadRequest)
			return
		}

		// генерируем временный url для отображения (будет обновлен после создания)
		photoURL = "/admin/masters/photo/temp"
	}

	if idStr != "" {
		id, convErr := strconv.Atoi(idStr)
		if convErr != nil {
			http.Error(w, "Неверный ID", http.StatusBadRequest)
			return
		}

		// если новая фотография не загружена, сохраняем старую
		existing, err := h.repo.GetByID(id)
		if err != nil {
			http.Error(w, "Мастер не найден", http.StatusNotFound)
			return
		}

		if len(photoData) == 0 {
			// сохраняем существующие данные
			photoData = existing.PhotoData
			photoType = existing.PhotoType
			photoURL = existing.PhotoURL
		}

		m := models.Master{
			ID:        id,
			FirstName: firstName,
			LastName:  lastName,
			PhotoURL:  photoURL,
			PhotoData: photoData,
			PhotoType: photoType,
		}

		if err := h.repo.Update(&m); err != nil {
			http.Error(w, "Ошибка обновления: "+err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		if len(photoData) == 0 {
			http.Error(w, "Фото обязательно", http.StatusBadRequest)
			return
		}

		m := models.Master{
			FirstName: firstName,
			LastName:  lastName,
			PhotoURL:  "/admin/masters/photo/temp",
			PhotoData: photoData,
			PhotoType: photoType,
		}

		if err := h.repo.Create(&m); err != nil {
			http.Error(w, "Ошибка создания: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// обновляем url фотографии с id (используем публичный маршрут)
		photoURL = fmt.Sprintf("/masters/photo/%d", m.ID)
		if err := h.repo.UpdatePhotoURL(m.ID, photoURL); err != nil {
			fmt.Printf("Ошибка обновления URL фотографии: %v\n", err)
		}
	}

	http.Redirect(w, r, "/admin/masters", http.StatusSeeOther)
}

func (h *AdminMastersHandler) GetPhoto(w http.ResponseWriter, r *http.Request) {
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
	w.Header().Set("Cache-Control", "public, max-age=31536000") // Кэширование на 1 год

	// отправляем данные изображения
	w.Write(master.PhotoData)
}

func (h *AdminMastersHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Ошибка парсинга формы", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}
	if err := h.repo.Delete(id); err != nil {
		http.Error(w, "Ошибка удаления: "+err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/admin/masters", http.StatusSeeOther)
}
