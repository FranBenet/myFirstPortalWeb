package main

import (
	"cars/pkg/routes"
	"log"
	"net/http"
)

func main() {

	if err := http.ListenAndServe(":8080", routes.Routes()); err != nil {
		log.Fatal(err)
	}

}
