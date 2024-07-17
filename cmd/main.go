package main

import (
	"cars/pkg/helpers"
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

	errChannel := make(chan error, 1)
	defer close(errChannel)

	go helpers.InitVariable(errChannel)
	err := <-errChannel

	if err != nil {
		fmt.Println("ERROR:", err)
	}
	fmt.Println("Running Server in 8080...")
	if err := http.ListenAndServe(":8080", routes.Routes()); err != nil {
		log.Fatal(err)
	}
}
