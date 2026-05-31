package middleware

import (
	"net/http"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// пока не научился пушить logger по middlware
		// fmt.Println(r.Method, r.URL.Path)

		next.ServeHTTP(w, r)
	})
}
