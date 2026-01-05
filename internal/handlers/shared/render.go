package shared

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/csrf"
)

var templateCache map[string]*template.Template

func SetTemplateCache(cache map[string]*template.Template) {
	templateCache = cache
}

func RenderTemplate(w http.ResponseWriter, r *http.Request, name string, data interface{}) {
	t, ok := templateCache[name] // шаблон из кэша
	if !ok {
		log.Printf("Шаблон %s не найден", name)
		http.Error(w, "Страница не найдена", http.StatusNotFound)
		return
	}

	// Prepare template data with CSRF token
	var templateData interface{}

	if data == nil {
		// If no data, create map with just CSRF token
		templateData = map[string]interface{}{
			"CSRFToken": csrf.Token(r),
			"CSRFField": template.HTML(csrf.TemplateField(r)),
		}
	} else {
		// If data exists, merge CSRF token into it
		switch v := data.(type) {
		case map[string]interface{}:
			// Already a map, add CSRF fields
			v["CSRFToken"] = csrf.Token(r)
			v["CSRFField"] = template.HTML(csrf.TemplateField(r))
			templateData = v
		default:
			// Struct or other type - wrap in map
			templateData = map[string]interface{}{
				"Data":      data,
				"CSRFToken": csrf.Token(r),
				"CSRFField": template.HTML(csrf.TemplateField(r)),
			}
		}
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := t.ExecuteTemplate(w, "base.html", templateData); err != nil {
		log.Printf("Ошибка рендеринга %s: %v", name, err)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
	}
}
