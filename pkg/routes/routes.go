package routes

import (
	"cars/pkg/handlers"
	"net/http"
)

func Routes() *http.ServeMux {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./web/static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", handlers.Homepage)
	// mux.HandleFunc("/gallery", handlers.Gallery)

	return mux
}
