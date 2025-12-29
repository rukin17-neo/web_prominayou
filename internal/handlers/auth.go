package handlers

import "net/http"

// TODO Доделать!!!
// mock auth
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// if r.Header.Get("X-Admin-Token") != "prom" {
		// 	http.Error(w, "Доступ запрещен", http.StatusForbidden)
		// 	return
		// }

		next.ServeHTTP(w, r)
	})
}
