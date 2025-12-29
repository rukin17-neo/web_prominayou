package handlers

import (
	"log"
	"net/http"
)

func RenderTemplate(w http.ResponseWriter, name string, data interface{}) {
	t, ok := templateCache[name] // шаблон из кэша
	if !ok {
		log.Printf("Шаблон %s не найден", name)
		http.Error(w, "Страница не найдена", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := t.ExecuteTemplate(w, "base.html", data); err != nil {
		log.Printf("Ошибка рендеринга %s: %v", name, err)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
	}
}
