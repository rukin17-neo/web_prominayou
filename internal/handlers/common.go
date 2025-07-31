package handlers

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

func RenderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	templateFiles := []string{
		filepath.Join("templates", "base.html"),
		filepath.Join("templates", tmpl),
	}

	t, err := template.ParseFiles(templateFiles...)
	if err != nil {
		log.Printf("Ошибка загрузки шаблонов: %v", err)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	err = t.ExecuteTemplate(w, "base.html", data)
	if err != nil {
		log.Printf("Ошибка рендеринга шаблона: %v", err)
		http.Error(w, "Ошибка отображения страницы", http.StatusInternalServerError)
	}
}
