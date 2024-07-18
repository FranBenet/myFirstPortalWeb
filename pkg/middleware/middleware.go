package middleware

import (
	"fmt"
	"net/http"
)

// ErrorHandler is middleware to handle errors and serve a custom error page
func ErrorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("Internal server error:", err)
				// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				// Serve a custom error page
				http.ServeFile(w, r, "./web/templates/500.html")
			}
		}()
		next.ServeHTTP(w, r)
	})
}
