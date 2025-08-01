package main

import (
	"log"
	"net/http"
	"prommsc/internal/handlers"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", handlers.HomeHandler).Methods("GET")
	r.HandleFunc("/contacts", handlers.ContactsHandler)
	r.HandleFunc("/price", handlers.PriceHandler)
	r.HandleFunc("/reviews", handlers.ReviewsHandler)

	r.NotFoundHandler = http.HandlerFunc(handlers.NotFoundHandler)

	r.PathPrefix("/static/").Handler(
		http.StripPrefix("/static", http.FileServer(http.Dir("static"))),
	)

	port := ":8006"
	log.Printf("Сервер запущен на http://localhost%s", port)
	log.Fatal(http.ListenAndServe(port, r)) //TODO ListenAndServeTLC для https
}
