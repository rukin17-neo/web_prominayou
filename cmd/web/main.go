package main

import (
	"log"
	"net/http"
	"prommsc/config"
	"prommsc/internal/handlers"
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
		log.Fatalf("ошибка подключения к БД: %v, err")
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Не удалось проверить подключение к БД: %v", err)
	}
	log.Println("Успешное подключение к PostgreSQL")

	serviceRepo := models.


	r := mux.NewRouter()

	r.HandleFunc("/", handlers.HomeHandler).Methods("GET")
	r.HandleFunc("/contacts", handlers.ContactsHandler)
	r.HandleFunc("/price", handlers.PriceHandler)
	r.HandleFunc("/reviews", handlers.ReviewsHandler)

	r.NotFoundHandler = http.HandlerFunc(handlers.NotFoundHandler)

	r.PathPrefix("/static/").Handler(
		http.StripPrefix("/static", http.FileServer(http.Dir("static"))),
	)

	fs := http.FileServer(http.Dir("static"))
	r.PathPrefix("/static/").Handler(
		http.StripPrefix("/static", cacheControl(fs)),
	)

	port := ":8006"
	log.Printf("Сервер запущен на http://localhost%s", port)
	log.Fatal(http.ListenAndServe(port, r)) //TODO ListenAndServeTLC для https
}
