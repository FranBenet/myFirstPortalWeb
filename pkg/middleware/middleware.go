package middleware

import (
	"cars/pkg/handlers"
	"fmt"
	"net/http"
)

func main() {

	// Wrap the mux with the errorHandler middleware
	wrappedRouter := middleware.ErrorHandler(router)
}

// ErrorHandler is middleware to handle errors and serve a custom error page
func ErrorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("Internal server error:", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				// Serve a custom error page
				handlers.NotFoundHandler(w, r)
				// http.ServeFile(w, r, "./web/templates/500.html")
				return
			}
		}()
		next.ServeHTTP(w, r)
	})
}
