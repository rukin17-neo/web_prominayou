package handlers

import (
	"html/template"
	"path/filepath"
	"prommsc/internal/handlers/shared"
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

		// Функции для шаблонов
		funcMap := template.FuncMap{
			"add": func(a, b int) int { return a + b },
			"sub": func(a, b int) int { return a - b },
			"seq": func(start, end int) []int {
				if start > end {
					return []int{}
				}
				seq := make([]int, end-start+1)
				for i := range seq {
					seq[i] = start + i
				}
				return seq
			},
		}

		basePath := filepath.Join("templates", "base.html")
		paginationPath := filepath.Join("templates", "partials", "pagination.html")

		// Парсим базовый шаблон с функциями
		baseTmpl, err := template.New("base.html").Funcs(funcMap).ParseFiles(basePath)
		if err != nil {
			initErr = err
			return
		}

		// Добавляем partial шаблон пагинации в базовый шаблон
		baseTmpl, err = baseTmpl.ParseFiles(paginationPath)
		if err != nil {
			initErr = err
			return
		}

		templates := []string{
			"index.html",
			"services.html",
			"contacts.html",
			"masters.html",
			"admin/admin.html",
			"admin/masters.html",
			"admin/create_service.html",
			"admin/edit_service.html",
			"admin/services.html",
			"admin/login.html",
			"admin/forgot_password.html",
			"admin/reset_password.html",
			"admin/users.html",
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

		// устанавливаем кэш в shared пакете
		shared.SetTemplateCache(templateCache)
	})
	return initErr
}
