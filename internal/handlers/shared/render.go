package shared

import (
	"html/template"
	"log"
	"net/http"
	"reflect"
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

	// Создаем map для данных
	dataMap := make(map[string]interface{})

	// Копируем все поля из оригинальной структуры
	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Struct {
		copyStructFields(val, dataMap)
	} else if val.Kind() == reflect.Map {
		// Если data уже map, копируем все поля
		iter := val.MapRange()
		for iter.Next() {
			k := iter.Key()
			v := iter.Value()
			dataMap[k.String()] = v.Interface()
		}
	}

	// CSRF removed - no CSRFField added

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := t.ExecuteTemplate(w, "base.html", dataMap); err != nil {
		log.Printf("Ошибка рендеринга %s: %v", name, err)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
	}
}

// copyStructFields рекурсивно копирует поля структуры, включая embedded fields
func copyStructFields(val reflect.Value, dataMap map[string]interface{}) {
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldVal := val.Field(i)

		if !field.IsExported() {
			continue
		}

		// Если поле является embedded struct, копируем его поля
		if field.Anonymous && fieldVal.Kind() == reflect.Struct {
			copyStructFields(fieldVal, dataMap)
		} else {
			dataMap[field.Name] = fieldVal.Interface()
		}
	}
}
