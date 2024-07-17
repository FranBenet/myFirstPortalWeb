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
	mux.HandleFunc("/id", handlers.SelectCar)
	mux.HandleFunc("/liked-compared", handlers.StatusChange)
	mux.HandleFunc("/comparePage", handlers.ComparePage)
	mux.HandleFunc("/lastCompare", handlers.LastCompare)
	mux.HandleFunc("/favouritePage", handlers.FavouritesPage)
	mux.HandleFunc("/search", handlers.Filter)

	return mux
}
