package admin

import (
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"prommsc/internal/handlers"
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

// Допустимые расширения и MIME типы для загружаемых изображений
var (
	allowedExtensions = map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".webp": true,
		".gif":  true,
	}

	allowedMimeTypes = map[string]string{
		"image/jpeg": "image/jpeg",
		"image/png":  "image/png",
		"image/webp": "image/webp",
		"image/gif":  "image/gif",
	}
)

// validateImageFile проверяет, что загруженный файл является валидным изображением
func validateImageFile(filename string, data []byte) (string, error) {
	// 1. Проверка расширения файла
	ext := strings.ToLower(filepath.Ext(filename))
	if !allowedExtensions[ext] {
		return "", fmt.Errorf("недопустимое расширение файла: %s. Разрешены только: .jpg, .jpeg, .png, .webp, .gif", ext)
	}

	// 2. Проверка реального содержимого файла по magic bytes
	// http.DetectContentType анализирует первые 512 байт для определения типа
	detectedType := http.DetectContentType(data)

	// 3. Проверка, что детектированный тип находится в белом списке
	mimeType, ok := allowedMimeTypes[detectedType]
	if !ok {
		return "", fmt.Errorf("недопустимый тип файла: %s. Файл не является изображением", detectedType)
	}

	// 4. Дополнительная проверка magic bytes для основных форматов
	if err := verifyImageMagicBytes(data, mimeType); err != nil {
		return "", err
	}

	return mimeType, nil
}

// verifyImageMagicBytes проверяет magic bytes (сигнатуры) файлов изображений
func verifyImageMagicBytes(data []byte, expectedType string) error {
	if len(data) < 12 {
		return fmt.Errorf("файл слишком маленький для проверки")
	}

	switch expectedType {
	case "image/jpeg":
		// JPEG начинается с FF D8 FF
		if !bytes.HasPrefix(data, []byte{0xFF, 0xD8, 0xFF}) {
			return fmt.Errorf("файл не является валидным JPEG изображением")
		}
	case "image/png":
		// PNG начинается с 89 50 4E 47 0D 0A 1A 0A
		if !bytes.HasPrefix(data, []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}) {
			return fmt.Errorf("файл не является валидным PNG изображением")
		}
	case "image/gif":
		// GIF начинается с GIF87a или GIF89a
		if !bytes.HasPrefix(data, []byte("GIF87a")) && !bytes.HasPrefix(data, []byte("GIF89a")) {
			return fmt.Errorf("файл не является валидным GIF изображением")
		}
	case "image/webp":
		// WebP: байты 0-3 = "RIFF", байты 8-11 = "WEBP"
		if len(data) < 12 || !bytes.HasPrefix(data, []byte("RIFF")) || !bytes.Equal(data[8:12], []byte("WEBP")) {
			return fmt.Errorf("файл не является валидным WebP изображением")
		}
	}

	return nil
}

func (h *AdminMastersHandler) List(w http.ResponseWriter, r *http.Request) {
	// Получение параметров пагинации из запроса
	paginationParams := models.NewPaginationParams(r)

	// Получение мастеров с пагинацией
	masters, pagination, err := h.repo.GetAllWithPagination(paginationParams)
	if err != nil {
		logAndRespondWithError(w, "GetAllMasters", err, ErrMsgLoadFailed, http.StatusInternalServerError)
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
		Title       string
		Masters     []models.Master
		EditMaster  *models.Master
		CurrentUser *models.User
		Pagination  models.PaginationResult
	}
	shared.RenderTemplate(w, r, "admin/masters.html", pageData{
		Title:       "Мастера",
		Masters:     masters,
		EditMaster:  editMaster,
		CurrentUser: handlers.GetCurrentUser(r),
		Pagination:  pagination,
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

		// валидация размера файла (максимум 5MB)
		if len(photoData) > 5*1024*1024 {
			http.Error(w, "Файл слишком большой. Максимальный размер: 5MB", http.StatusBadRequest)
			return
		}

		// безопасная валидация файла:
		// - проверка расширения из белого списка
		// - проверка реального содержимого по magic bytes
		// - защита от подделки Content-Type заголовка
		validatedType, err := validateImageFile(header.Filename, photoData)
		if err != nil {
			http.Error(w, "Ошибка валидации файла: "+err.Error(), http.StatusBadRequest)
			return
		}

		photoType = validatedType

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
			photoType = existing.PhotoType.String
			photoURL = existing.PhotoURL.String
		}

		m := models.Master{
			ID:        id,
			FirstName: firstName,
			LastName:  lastName,
			PhotoURL:  sql.NullString{String: photoURL, Valid: photoURL != ""},
			PhotoData: photoData,
			PhotoType: sql.NullString{String: photoType, Valid: photoType != ""},
		}

		if err := h.repo.Update(&m); err != nil {
			logAndRespondWithError(w, "UpdateMaster", err, ErrMsgUpdateFailed, http.StatusInternalServerError)
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
			PhotoURL:  sql.NullString{String: "/admin/masters/photo/temp", Valid: true},
			PhotoData: photoData,
			PhotoType: sql.NullString{String: photoType, Valid: true},
		}

		if err := h.repo.Create(&m); err != nil {
			logAndRespondWithError(w, "CreateMaster", err, ErrMsgCreateFailed, http.StatusInternalServerError)
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
	shared.ServePhoto(h.repo)(w, r)
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
		logAndRespondWithError(w, "DeleteMaster", err, ErrMsgDeleteFailed, http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/admin/masters", http.StatusSeeOther)
}
