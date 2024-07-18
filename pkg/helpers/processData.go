package helpers

import (
	"cars/pkg/config"
	"cars/pkg/models"
	"fmt"
	"log"
	"net/http"
	"strings"
	"text/template"
)

// Takes one variable type models.Car (which has the same structure as the API)
// and returns a variable type models.Card with all the information needed.
func CreateSmallCard(car models.Car) (models.Card, error) {

	manufacturerChannel := make(chan models.Manufacturers, 1)
	manufacturerErrChannel := make(chan error, 1)
	categoryChannel := make(chan models.Categories, 1)
	categoryErrChannel := make(chan error, 1)

	go FetchManufacturer(car.ManufacturerID, manufacturerChannel, manufacturerErrChannel)
	go FetchCategory(car.CategoryID, categoryChannel, categoryErrChannel)

	manufacturer := <-manufacturerChannel
	err := <-manufacturerErrChannel
	if err != nil {
		fmt.Println("Error fetching data from the API")
		return models.Card{}, err
	}

	category := <-categoryChannel
	err = <-categoryErrChannel
	if err != nil {
		fmt.Println("Error fetching data from the API")
		return models.Card{}, err
	}

	var card models.Card

	card.Id = car.Id
	card.Name = car.Name
	card.Year = car.Year
	card.Image = car.Image
	card.Category = category.Name
	card.Manufacturer = manufacturer.Name

	//	Liked and Compared are boolean values that will allow the HTML to determine
	//	the appearance for the correspondent icons.
	//	We get the values from the maps initialized when running the application.
	card.Liked = config.FavouritesMap[car.Id]
	card.Compared = config.ComparisonMap[car.Id]

	return card, nil
}

func CreateSmallCardsBatch(carsSelected []models.Car) ([]models.Card, error) {
	var cards []models.Card
	for _, car := range carsSelected {
		card, err := CreateSmallCard(car)
		if err != nil {
			fmt.Printf("Error creating small card: %v", err)
			return []models.Card{}, err
		}
		cards = append(cards, card)
	}
	return cards, nil
}

// Takes one variable type models.Car (which has the same structure as the API)
// and returns a variable type models.ExtendedCard with all the extended information wanted.
func CreateBigCard(car models.Car) (models.ExtendedCard, error) {

	manufacturerChannel := make(chan models.Manufacturers, 1)
	errManufacturerChannel := make(chan error, 1)
	categoryChannel := make(chan models.Categories, 1)
	errCategoryChannel := make(chan error, 1)

	go FetchCategory(car.CategoryID, categoryChannel, errCategoryChannel)
	go FetchManufacturer(car.ManufacturerID, manufacturerChannel, errManufacturerChannel)

	category := <-categoryChannel
	err := <-errCategoryChannel
	if err != nil {
		fmt.Println("Error fetching data from the API")
		return models.ExtendedCard{}, err
	}

	manufacturer := <-manufacturerChannel
	err = <-errManufacturerChannel

	if err != nil {
		fmt.Println("Error fetching data from the API")
		return models.ExtendedCard{}, err
	}

	var card models.ExtendedCard
	card.Id = car.Id
	card.Name = car.Name
	card.Year = car.Year
	card.Image = car.Image
	card.Category = category.Name
	card.Manufacturer = manufacturer.Name
	card.FoundingYear = manufacturer.FoundingYear
	card.Country = manufacturer.Country
	card.Engine = car.Specifications.Engine
	card.Horsepower = car.Specifications.Horsepower
	card.Transmission = car.Specifications.Transmission
	card.DriveTrain = car.Specifications.DriveTrain

	//	Liked and Compared are boolean values that will allow the HTML to determine
	//	the appearance for the correspondent icons.
	//	We get their values from the maps initialized when running the application.
	card.Liked = config.FavouritesMap[car.Id]
	card.Compared = config.ComparisonMap[car.Id]

	return card, nil
}

func CreateBigCardsBatch(carsSelected []models.Car) ([]models.ExtendedCard, error) {
	var cards []models.ExtendedCard
	for _, car := range carsSelected {
		card, err := CreateBigCard(car)
		if err != nil {
			log.Printf("Error creating big card: %v", err)
			return []models.ExtendedCard{}, err
		}
		cards = append(cards, card)
	}
	return cards, nil
}

// Initializes the global variables FavouritesMap, ComparisonMap, CategoriesFilterMap, ManufacturersFilterMap and ModelsFilterMap.
func InitVariable(errChannel chan error) {
	carsDataChannel := make(chan []models.Car, 1)
	carsErrChannel := make(chan error, 1)

	manufacturersChannel := make(chan []models.Manufacturers, 1)
	manufacturersErrChannel := make(chan error, 1)

	categoriesChannel := make(chan []models.Categories, 1)
	categoriesErrChannel := make(chan error, 1)

	var err error

	go FetchCars(carsDataChannel, carsErrChannel)
	go FetchManufacturers(manufacturersChannel, manufacturersErrChannel)
	go FetchCategories(categoriesChannel, categoriesErrChannel)

	carsData := <-carsDataChannel
	err = <-carsErrChannel
	if err != nil {
		fmt.Println("Error fetching data from the API")
		errChannel <- err
		return
	}

	manufacturersData := <-manufacturersChannel
	err = <-manufacturersErrChannel
	if err != nil {
		fmt.Println("Error fetching data from the API")
		errChannel <- err
		return
	}

	categoriesData := <-categoriesChannel
	err = <-categoriesErrChannel
	if err != nil {
		fmt.Println("Error fetching data from the API")
		errChannel <- err
		return
	}

	go func() {
		for i, car := range carsData {
			config.FavouritesMap[car.Id] = false
			config.ComparisonMap[car.Id] = false
			config.ModelsFilterMap[car.Name] = false
			config.TotalNumCars = i
		}
	}()

	go func() {
		for _, manufacturer := range manufacturersData {
			config.ManufacturersFilterMap[manufacturer.Id] = false
		}
	}()

	go func() {
		for _, category := range categoriesData {
			config.CategoriesFilterMap[category.Id] = false
		}
	}()

	errChannel <- nil
	close(errChannel)
}

func SearchQueryCars(query string) ([]models.Car, error) {

	//	Fetch all the cars.
	carsDataChannel := make(chan []models.Car, 1)
	errChannel := make(chan error, 1)

	FetchCars(carsDataChannel, errChannel)

	cars := <-carsDataChannel
	err := <-errChannel
	if err != nil {
		fmt.Println("Error fetching data from the API")
		return []models.Car{}, err
	}

	//	Collect all the cars that matches the query
	var filteredCars []models.Car
	query = strings.ToLower(query)
	for _, car := range cars {
		manufacturerChannel := make(chan models.Manufacturers, len(cars))
		errManufactureChannel := make(chan error, len(cars))
		categoryChannel := make(chan models.Categories, len(cars))
		errCategoryChannel := make(chan error, len(cars))

		go FetchManufacturer(car.ManufacturerID, manufacturerChannel, errManufactureChannel)
		go FetchCategory(car.CategoryID, categoryChannel, errCategoryChannel)

		manufacturer := <-manufacturerChannel
		err := <-errManufactureChannel
		if err != nil {
			fmt.Println("Error fetching data from the API")
			return []models.Car{}, err
		}

		category := <-categoryChannel
		err = <-errCategoryChannel
		if err != nil {
			fmt.Println("Error fetching data from the API")
			return []models.Car{}, err
		}

		if strings.Contains(strings.ToLower(car.Name), query) || strings.Contains(strings.ToLower(manufacturer.Name), query) || strings.Contains(strings.ToLower(category.Name), query) {
			filteredCars = append(filteredCars, car)
		}
	}
	return filteredCars, nil
}

func RenderTemplate(w http.ResponseWriter, htmlTemplate []string, name string, data models.DataResponse) {

	tmpl, err := template.ParseFiles(htmlTemplate...)
	if err != nil {
		fmt.Printf("Error Parsing Template: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = tmpl.ExecuteTemplate(w, name, data)
	if err != nil {
		fmt.Printf("Error Executing Template: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
