package handlers

import (
	"html/template"
	"path/filepath"
	"sync"
)

var (
	templateCache     map[string]*template.Template // кэш
	templateCacheOnce sync.Once
)

func InitTemplates() error {
	var initErr error
	templateCacheOnce.Do(func() {
		templateCache = make(map[string]*template.Template)

		basePath := filepath.Join("templates", "base.html")
		baseTmpl := template.Must(template.ParseFiles(basePath))

		templates := []string{
			"index.html",
			"services.html",
			"contacts.html",
			"reviews.html",
			"admin/create_service.html",
			"admin/edit_service.html",
			"admin/services.html",
		}

		for _, tmpl := range templates {
			path := filepath.Join("templates", tmpl)

			// клонируем базовый шаблон и добавляем контент
			t, err := baseTmpl.Clone()
			if err != nil {
				initErr = err
				return
			}

			t, err = t.ParseFiles(path)
			if err != nil {
				initErr = err
				return
			}

			templateCache[tmpl] = t
		}
	})
	return initErr
}
