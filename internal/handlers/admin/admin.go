package admin

import (
	"net/http"
	"prommsc/internal/handlers/shared"
	"prommsc/models"
)

func AdminDashboard(w http.ResponseWriter, r *http.Request) {
	data := models.PageData{
		Title:   "Панель управления",
		Content: "Выберите раздел для управления",
	}
	shared.RenderTemplate(w, "admin/admin.html", data)
}
