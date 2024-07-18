package main

import (
	"cars/pkg/helpers"
	"cars/pkg/middleware"
	"cars/pkg/routes"
	"fmt"
	"log"
	"net/http"
)

func main() {
	//	We populate the variables FavouritesMap and ComparisonMap with
	//	a list of all the ID cars and a boolean value initiated as false.
	//	This maps will help us keep track of the items liked, and items selected to
	//	be compared.

	// Get the router from the routes package
	router := routes.Routes()

	// Wrap the mux with the errorHandler middleware
	wrappedRouter := middleware.ErrorHandler(router)

	errChannel := make(chan error, 1)
	defer close(errChannel)

	go helpers.InitVariable(errChannel)
	err := <-errChannel
	if err != nil {
		fmt.Println("Error initiating program.")
		log.Fatal(err)
	}

	fmt.Println("Running Server in 8080...")
	if err := http.ListenAndServe(":8080", wrappedRouter); err != nil {
		log.Fatal(err)
	}
}
