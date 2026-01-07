package main

import (
	"log"
	"net/http"
	"prommsc/config"
	"prommsc/internal/auth"
	"prommsc/internal/handlers"
	adminHandlers "prommsc/internal/handlers/admin"
	clientHandlers "prommsc/internal/handlers/client"
	"prommsc/internal/middleware"
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

	// Initialize rate limiters
	if err := middleware.InitRateLimiters(); err != nil {
		log.Fatalf("Failed to initialize rate limiters: %v", err)
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

	// Применяем security headers ко всем маршрутам
	r.Use(middleware.SecurityHeaders())

	r.HandleFunc("/", clientHandlers.HomeHandler).Methods("GET")
	r.HandleFunc("/services", servicesHandler.GetAllServices).Methods("GET")
	r.HandleFunc("/contacts", clientHandlers.ContactsHandler).Methods("GET")
	r.HandleFunc("/reviews", clientHandlers.ReviewsHandler).Methods("GET")
	r.HandleFunc("/masters", mastersHandler.List).Methods("GET")

	r.HandleFunc("/masters/photo/{id}", mastersHandler.GetPhoto).Methods("GET")

	// Admin subrouter
	adminRouter := r.PathPrefix("/admin").Subrouter()
	// CSRF Protection removed

	// Public auth routes on admin subrouter (get CSRF protection)
	// Login with strict rate limiting
	adminRouter.HandleFunc("/login", authHandler.LoginPage).Methods("GET")
	adminRouter.Handle("/login", middleware.RateLimitLogin()(
		http.HandlerFunc(authHandler.Login))).Methods("POST")

	// Forgot password with strict rate limiting
	adminRouter.HandleFunc("/forgot-password", authHandler.ForgotPasswordPage).Methods("GET")
	adminRouter.Handle("/forgot-password", middleware.RateLimitForgotPassword()(
		http.HandlerFunc(authHandler.ForgotPassword))).Methods("POST")

	// Reset password with strict rate limiting
	adminRouter.HandleFunc("/reset-password", authHandler.ResetPasswordPage).Methods("GET")
	adminRouter.Handle("/reset-password", middleware.RateLimitResetPassword()(
		http.HandlerFunc(authHandler.ResetPassword))).Methods("POST")

	// Protected routes - create nested subrouter with AuthMiddleware
	protectedRouter := adminRouter.NewRoute().Subrouter()
	protectedRouter.Use(handlers.AuthMiddleware)

	// Dashboard
	protectedRouter.HandleFunc("", adminDashboard).Methods("GET")

	// Masters with CRUD rate limiting
	adminMasters := adminHandlers.NewAdminMastersHandler(mastersRepo)
	protectedRouter.Handle("/masters", middleware.RateLimitCRUD()(
		http.HandlerFunc(adminMasters.List))).Methods("GET")
	protectedRouter.Handle("/masters", middleware.RateLimitCRUD()(
		http.HandlerFunc(adminMasters.CreateOrUpdate))).Methods("POST")
	protectedRouter.Handle("/masters/delete", middleware.RateLimitCRUD()(
		http.HandlerFunc(adminMasters.Delete))).Methods("POST")
	protectedRouter.HandleFunc("/masters/photo/{id}", adminMasters.GetPhoto).Methods("GET")

	// Services with CRUD rate limiting
	protectedRouter.Handle("/services", middleware.RateLimitCRUD()(
		http.HandlerFunc(adminHandler.ListServices))).Methods("GET")
	protectedRouter.HandleFunc("/services/new", adminHandler.CreateServiceForm).Methods("GET")
	protectedRouter.Handle("/services", middleware.RateLimitCRUD()(
		http.HandlerFunc(adminHandler.CreateService))).Methods("POST")
	protectedRouter.HandleFunc("/services/edit/{id}", adminHandler.UpdateServiceForm).Methods("GET")
	protectedRouter.Handle("/services/{id}", middleware.RateLimitCRUD()(
		http.HandlerFunc(adminHandler.UpdateService))).Methods("POST", "PUT")
	protectedRouter.Handle("/services/delete/{id}", middleware.RateLimitCRUD()(
		http.HandlerFunc(adminHandler.DeleteService))).Methods("POST")

	// Logout (no rate limiting needed)
	protectedRouter.HandleFunc("/logout", authHandler.Logout).Methods("POST")

	// Users with high priority rate limiting
	protectedRouter.Handle("/users", middleware.RateLimitUserManagement()(
		http.HandlerFunc(usersHandler.List))).Methods("GET")
	protectedRouter.Handle("/users", middleware.RateLimitUserManagement()(
		http.HandlerFunc(usersHandler.CreateOrUpdate))).Methods("POST")
	protectedRouter.Handle("/users/delete", middleware.RateLimitUserManagement()(
		http.HandlerFunc(usersHandler.Delete))).Methods("POST")

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

	port := ":8004"
	log.Printf("Сервер запущен на http://localhost%s", port)
	log.Fatal(http.ListenAndServe(port, r))
}
