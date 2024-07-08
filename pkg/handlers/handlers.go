package handlers

import (
	"cars/pkg/models"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type data struct {
	Card        []models.CarModels
	Manufacture []string
	Category    []string
}

func Homepage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
		return
	}

	carModel, err := http.Get("http://localhost:3000/api/models")
	if err != nil {
		fmt.Println("Error getting carModels information from the cars API")
	}
	defer carModel.Body.Close()

	var carsData []models.CarModels

	if err = json.NewDecoder(carModel.Body).Decode(&carsData); err != nil {
		log.Fatalf("Failed to decode JSON carsData: %v", err)
	}

	category, err := http.Get("http://localhost:3000/api/categories")
	if err != nil {
		fmt.Println("Error getting car categories information from the cars API")
	}
	defer category.Body.Close()

	var categories []models.Categories

	if err = json.NewDecoder(category.Body).Decode(&categories); err != nil {
		log.Fatalf("Failed to decode JSON categories: %v", err)
	}

	var categoriesData []string
	for _, item := range categories {
		categoriesData = append(categoriesData, item.Name)
	}

	fmt.Println(categoriesData)

	brands, err := http.Get("http://localhost:3000/api/manufacturers")
	if err != nil {
		fmt.Println("Error getting car categories information from the cars API")
	}
	defer brands.Body.Close()

	var manufacturers []models.Manufacturers

	if err = json.NewDecoder(brands.Body).Decode(&manufacturers); err != nil {
		log.Fatalf("Failed to decode JSON categories: %v", err)
	}

	var manufacturersData []string

	for _, item := range manufacturers {
		manufacturersData = append(manufacturersData, item.Name)
	}

	fmt.Println(manufacturersData)

	var finalData data
	finalData.Card = carsData
	finalData.Category = categoriesData
	finalData.Manufacture = manufacturersData

	fmt.Println(finalData)
	fmt.Println()

	htmlTemplates := []string{
		"web/templates/index.html",
		"web/templates/card-template.html",
		"web/templates/main-bar.html",
		"web/templates/filter.html",
	}

	tmpl, err := template.ParseFiles(htmlTemplates...)
	if err != nil {
		log.Fatal(err)
	}

	err = tmpl.ExecuteTemplate(w, "index.html", finalData)
	if err != nil {
		fmt.Println("Oops! Something failed when Executing template")
		http.Error(w, "Internal Server Error Fran", http.StatusInternalServerError)
		return
	}
}

// func Gallery(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("This is GALLERY")
// }
