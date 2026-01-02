package main

import (
	"log"
	"net/http"
	"prommsc/config"
	"prommsc/internal/auth"
	"prommsc/internal/handlers"
	adminHandlers "prommsc/internal/handlers/admin"
	clientHandlers "prommsc/internal/handlers/client"
	"prommsc/models"
	"time"

	"github.com/gorilla/mux"
)

func cacheControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Кеширование на 1 год для статики
		w.Header().Set("ETag", `"`+time.Now().Format("20060102")+`"`)
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		w.Header().Set("Vary", "Accept-Encoding")
		h.ServeHTTP(w, r)
	})
}

func main() {
	config.LoadEnv()

	db, err := config.InitDB()
	if err != nil {
		log.Fatalf("ошибка подключения к БД: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Не удалось проверить подключение к БД: %v", err)
	}
	log.Println("Успешное подключение к PostgreSQL")

	if err := handlers.InitTemplates(); err != nil {
		log.Fatalf("Ошибка загрузки шаблонов: %v", err)
	}

	serviceRepo := models.NewServiceRepository(db)
	servicesHandler := clientHandlers.NewServicesHandler(serviceRepo)
	adminHandler := adminHandlers.NewAdminServicesHandler(serviceRepo)

	mastersRepo := models.NewMastersRepository(db)
	if err := mastersRepo.InitSchema(); err != nil {
		log.Fatalf("Не удалось создать таблицу masters: %v", err)
	}
	mastersHandler := clientHandlers.NewMastersHandler(mastersRepo)

	// Инициализация auth repositories
	userRepo := models.NewUserRepository(db)
	if err := userRepo.InitSchema(); err != nil {
		log.Fatalf("Не удалось создать таблицу users: %v", err)
	}

	sessionRepo := models.NewSessionRepository(db)
	if err := sessionRepo.InitSchema(); err != nil {
		log.Fatalf("Не удалось создать таблицу sessions: %v", err)
	}

	// Создание session manager
	sessionManager := auth.NewSessionManager(sessionRepo, userRepo)
	handlers.SetSessionManager(sessionManager)

	// Создание auth handlers
	authHandler := adminHandlers.NewAdminAuthHandler(userRepo, sessionManager)
	usersHandler := adminHandlers.NewAdminUsersHandler(userRepo)

	// admin handler
	adminDashboard := adminHandlers.AdminDashboard

	r := mux.NewRouter()

	r.HandleFunc("/", clientHandlers.HomeHandler).Methods("GET")
	r.HandleFunc("/services", servicesHandler.GetAllServices).Methods("GET")
	r.HandleFunc("/contacts", clientHandlers.ContactsHandler).Methods("GET")
	r.HandleFunc("/reviews", clientHandlers.ReviewsHandler).Methods("GET")
	r.HandleFunc("/masters", mastersHandler.List).Methods("GET")

	r.HandleFunc("/masters/photo/{id}", mastersHandler.GetPhoto).Methods("GET")

	// Публичные auth маршруты (ДО middleware)
	r.HandleFunc("/admin/login", authHandler.LoginPage).Methods("GET")
	r.HandleFunc("/admin/login", authHandler.Login).Methods("POST")
	r.HandleFunc("/admin/forgot-password", authHandler.ForgotPasswordPage).Methods("GET")
	r.HandleFunc("/admin/forgot-password", authHandler.ForgotPassword).Methods("POST")
	r.HandleFunc("/admin/reset-password", authHandler.ResetPasswordPage).Methods("GET")
	r.HandleFunc("/admin/reset-password", authHandler.ResetPassword).Methods("POST")

	adminRouter := r.PathPrefix("/admin").Subrouter()
	adminRouter.Use(handlers.AuthMiddleware)

	// /admin
	adminRouter.HandleFunc("", adminDashboard).Methods("GET")

	// /admin/masters
	adminMasters := adminHandlers.NewAdminMastersHandler(mastersRepo)
	adminRouter.HandleFunc("/masters", adminMasters.List).Methods("GET")
	adminRouter.HandleFunc("/masters", adminMasters.CreateOrUpdate).Methods("POST")
	adminRouter.HandleFunc("/masters/delete", adminMasters.Delete).Methods("POST")

	adminRouter.HandleFunc("/masters/photo/{id}", adminMasters.GetPhoto).Methods("GET")

	// /admin/services
	adminRouter.HandleFunc("/services", adminHandler.ListServices).Methods("GET")
	adminRouter.HandleFunc("/services/new", adminHandler.CreateServiceForm).Methods("GET")
	adminRouter.HandleFunc("/services", adminHandler.CreateService).Methods("POST")
	adminRouter.HandleFunc("/services/edit/{id}", adminHandler.UpdateServiceForm).Methods("GET")
	adminRouter.HandleFunc("/services/{id}", adminHandler.UpdateService).Methods("POST", "PUT")
	adminRouter.HandleFunc("/services/delete/{id}", adminHandler.DeleteService).Methods("POST")

	// /admin/logout (защищённый)
	adminRouter.HandleFunc("/logout", authHandler.Logout).Methods("POST")

	// /admin/users (защищённые)
	adminRouter.HandleFunc("/users", usersHandler.List).Methods("GET")
	adminRouter.HandleFunc("/users", usersHandler.CreateOrUpdate).Methods("POST")
	adminRouter.HandleFunc("/users/delete", usersHandler.Delete).Methods("POST")

	r.PathPrefix("/static/").Handler(
		http.StripPrefix("/static", cacheControl(http.FileServer(http.Dir("static")))),
	)

	// Фоновая очистка истекших сессий (каждый час)
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			if err := sessionManager.CleanupExpiredSessions(); err != nil {
				log.Printf("Ошибка очистки сессий: %v", err)
			} else {
				log.Println("Очистка истекших сессий выполнена")
			}
		}
	}()

	port := ":8003"
	log.Printf("Сервер запущен на http://localhost%s", port)
	log.Fatal(http.ListenAndServe(port, r))
}
