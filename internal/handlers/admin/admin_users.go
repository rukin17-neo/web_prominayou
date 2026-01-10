package admin

import (
	"net/http"
	"prommsc/internal/auth"
	"prommsc/internal/handlers"
	"prommsc/internal/handlers/shared"
	"prommsc/models"
	"strconv"
	"strings"
)

type AdminUsersHandler struct {
	userRepo *models.UserRepository
}

func NewAdminUsersHandler(userRepo *models.UserRepository) *AdminUsersHandler {
	return &AdminUsersHandler{userRepo: userRepo}
}

// List отображает список пользователей с возможностью редактирования
func (h *AdminUsersHandler) List(w http.ResponseWriter, r *http.Request) {
	// Получение параметров пагинации из запроса
	paginationParams := models.NewPaginationParams(r)

	// Получение пользователей с пагинацией
	users, pagination, err := h.userRepo.GetAllWithPagination(paginationParams)
	if err != nil {
		logAndRespondWithError(w, "GetAllUsers", err, ErrMsgLoadFailed, http.StatusInternalServerError)
		return
	}

	var editUser *models.User
	// Если есть параметр ?edit={id} - получаем пользователя для редактирования
	if editID := r.URL.Query().Get("edit"); editID != "" {
		if id, err := strconv.Atoi(editID); err == nil {
			if u, e := h.userRepo.GetByID(id); e == nil {
				editUser = u
			}
		}
	}

	// Получение текущего пользователя из контекста
	currentUser := handlers.GetCurrentUser(r)
	var currentUserID int
	if currentUser != nil {
		currentUserID = currentUser.ID
	}

	type pageData struct {
		Title         string
		Users         []models.User
		EditUser      *models.User
		CurrentUserID int
		CurrentUser   *models.User
		Pagination    models.PaginationResult
	}

	shared.RenderTemplate(w, r, "admin/users.html", pageData{
		Title:         "Пользователи",
		Users:         users,
		EditUser:      editUser,
		CurrentUserID: currentUserID,
		CurrentUser:   currentUser,
		Pagination:    pagination,
	})
}

// CreateOrUpdate создаёт нового пользователя или обновляет существующего
func (h *AdminUsersHandler) CreateOrUpdate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Ошибка парсинга формы", http.StatusBadRequest)
		return
	}

	username := strings.TrimSpace(r.FormValue("username"))
	email := strings.TrimSpace(r.FormValue("email"))
	password := r.FormValue("password")
	idStr := r.FormValue("id")

	// Базовая валидация
	if username == "" || email == "" {
		http.Error(w, "Имя пользователя и email обязательны", http.StatusBadRequest)
		return
	}

	// Валидация email (regex)
	if !isValidEmail(email) {
		http.Error(w, "Неверный формат email", http.StatusBadRequest)
		return
	}

	// Режим редактирования
	if idStr != "" {
		id, convErr := strconv.Atoi(idStr)
		if convErr != nil {
			http.Error(w, "Неверный ID", http.StatusBadRequest)
			return
		}

		// Получение существующего пользователя
		existingUser, err := h.userRepo.GetByID(id)
		if err != nil {
			http.Error(w, "Пользователь не найден", http.StatusNotFound)
			return
		}

		// Обновление данных
		existingUser.Username = username
		existingUser.Email = email

		if err := h.userRepo.Update(existingUser); err != nil {
			// Проверка на duplicate username/email
			if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
				http.Error(w, "Имя пользователя или email уже используется", http.StatusBadRequest)
				return
			}
			logAndRespondWithError(w, "UpdateUser", err, ErrMsgUpdateFailed, http.StatusInternalServerError)
			return
		}

		// Обновление пароля (если указан)
		if password != "" {
			// Валидация пароля (длина, сложность)
			if err := validatePassword(password); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			passwordHash, err := auth.HashPassword(password)
			if err != nil {
				http.Error(w, "Ошибка хэширования пароля", http.StatusInternalServerError)
				return
			}

			if err := h.userRepo.UpdatePassword(id, passwordHash); err != nil {
				logAndRespondWithError(w, "UpdatePassword", err, ErrMsgUpdateFailed, http.StatusInternalServerError)
				return
			}
		}

		http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
		return
	}

	// Режим создания
	if password == "" {
		http.Error(w, "Пароль обязателен при создании пользователя", http.StatusBadRequest)
		return
	}

	// Валидация пароля (длина, сложность)
	if err := validatePassword(password); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Хэширование пароля
	passwordHash, err := auth.HashPassword(password)
	if err != nil {
		http.Error(w, "Ошибка хэширования пароля", http.StatusInternalServerError)
		return
	}

	// Создание пользователя
	user := &models.User{
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
	}

	if err := h.userRepo.Create(user); err != nil {
		// Проверка на duplicate username/email
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
			http.Error(w, "Имя пользователя или email уже используется", http.StatusBadRequest)
			return
		}
		logAndRespondWithError(w, "CreateUser", err, ErrMsgCreateFailed, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
}

// Delete удаляет пользователя
func (h *AdminUsersHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Ошибка парсинга формы", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	// Получение текущего пользователя
	currentUser := handlers.GetCurrentUser(r)
	if currentUser != nil && currentUser.ID == id {
		http.Error(w, "Нельзя удалить текущего пользователя", http.StatusBadRequest)
		return
	}

	// Удаление пользователя
	if err := h.userRepo.Delete(id); err != nil {
		logAndRespondWithError(w, "DeleteUser", err, ErrMsgDeleteFailed, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
}
