package handlers

import (
	"net/http"
	"prommsc/models"
)

func PriceHandler(w http.ResponseWriter, r *http.Request) {
	service := models.Service{
		Name:     "Антицеллюлитный массаж",
		Price:    4000.0,
		Duration: "60 минут",
	}
	data := struct {
		models.PageData
		Service models.Service
	}{
		PageData: models.PageData{Title: "Услуга"},
		Service:  service,
	}
	RenderTemplate(w, "price.html", data)
}
