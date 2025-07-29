package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"prommsc/models"

	"github.com/gorilla/mux"
)

type PageData struct {
	Title   string
	Content string
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	templateFiles := []string{
		filepath.Join("templates", "base.html"),
		filepath.Join("templates", tmpl),
	}

	log.Printf("Loading templates: %v", templateFiles)

	t, err := template.ParseFiles(templateFiles...)
	if err != nil {
		http.Error(w, "Ошибка загрузки шаблонов: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = t.ExecuteTemplate(w, "base.html", data)
	if err != nil {
		http.Error(w, "Ошибка выполнения шаблона: "+err.Error(), http.StatusInternalServerError)
	}
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index.html", PageData{
		Title:   "ProminaYou - Главная",
		Content: "Добро пожаловать на наш сайт!",
	})
}

func ContactsHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "contacts.html", PageData{
		Title:   "Контакты",
		Content: "+7495 747 4747",
	})
}

func PriceHandler(w http.ResponseWriter, r *http.Request) {
	service := models.Service{
		Name:  "Антицеллюлитный массаж",
		Price: "4000 руб.",
		Time:  "60 минут",
	}
	data := struct {
		PageData
		Service models.Service
	}{
		PageData: PageData{Title: "Услуга"},
		Service:  service,
	}
	renderTemplate(w, "price.html", data)
}

func ReviewsHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "reviews.html", PageData{
		Title:   "Отзывы",
		Content: "Отлично!",
	})
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler).Methods("GET")
	r.HandleFunc("/contacts", ContactsHandler)
	r.HandleFunc("/price", PriceHandler)
	r.HandleFunc("/reviews", ReviewsHandler)

	r.PathPrefix("/static/").Handler(
		http.StripPrefix("/static", http.FileServer(http.Dir("static"))),
	)

	http.Handle("/", r)

	port := ":8002"
	log.Printf("Сервер запущен на http://localhost%s", port)
	log.Fatal(http.ListenAndServe(port, r)) // TODO ListenAndServeTLS - для https!!!
}
