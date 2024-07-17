package helpers

import (
	"cars/pkg/config"
	"cars/pkg/models"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
)

// Fetch a car from the API by ID.
func FetchCar(id int, carChannel chan models.Car, errChannel chan error) {
	// var mu sync.Mutex
	// mu.Lock()
	// defer mu.Unlock()
	//	Convert the ID from int to string.
	idString := strconv.Itoa(id)

	// Attach the ID in string format to the URL to fetch.
	car, err := http.Get("http://localhost:3000/api/models/" + idString)
	if err != nil {
		log.Printf("Error getting car categories information from the cars API: %v", err)
		carChannel <- models.Car{}
		errChannel <- err
		// close(carChannel)
		// close(errChannel)
		return
	}
	defer car.Body.Close()

	//	Write the data in []byte type.
	data, err := io.ReadAll(car.Body)
	if err != nil {
		log.Printf("Error. FetchCar() - io.ReadAll error: %v", err)
		carChannel <- models.Car{}
		errChannel <- err
		// close(carChannel)
		// close(errChannel)
		return
	}

	var carData models.Car
	//	Decode the json data and store it in a variable with same fields as the json data.

	err = json.Unmarshal(data, &carData)
	if err != nil {
		log.Printf("Error. FetchCar() - json.Unmarshall error: %v", err)
		carChannel <- models.Car{}
		errChannel <- err
		// close(carChannel)
		// close(errChannel)
		return
	}

	carChannel <- carData
	errChannel <- nil
	// close(carChannel)
	// close(errChannel)
	fmt.Println("Car Fetched" + carData.Name)
}

// Fetch all cars from the API.
func FetchCars(carsDataChannel chan []models.Car, errChannel chan error) {

	fmt.Println("Starting to Fetch Cars...")
	cars, err := http.Get("http://localhost:3000/api/models")
	if err != nil {
		log.Printf("Error getting carModels information from the cars API: %v", err)
		carsDataChannel <- nil
		errChannel <- err
		// close(carsDataChannel)
		// close(errChannel)
		return
	}
	defer cars.Body.Close()

	var carsData []models.Car

	if err = json.NewDecoder(cars.Body).Decode(&carsData); err != nil {
		log.Fatalf("Failed to decode JSON carsData: %v", err)
		carsDataChannel <- nil
		errChannel <- err
		// close(carsDataChannel)
		// close(errChannel)
		return
	}

	carsDataChannel <- carsData
	errChannel <- nil
	// close(carsDataChannel)
	// close(errChannel)
	fmt.Println("Cars Fetched.")
}

// Fetch a category from the API by ID.
func FetchCategory(id int, categoryChannel chan models.Categories, errChannel chan error) {
	fmt.Printf("Starting to Fetch Category %d...", id)
	idString := strconv.Itoa(id)

	category, err := http.Get("http://localhost:3000/api/categories/" + idString)
	if err != nil {
		log.Printf("Error getting car categories information from the cars API: %v", err)
		categoryChannel <- models.Categories{}
		errChannel <- err
		// close(categoryChannel)
		// close(errChannel)
		return
	}
	defer category.Body.Close()

	var categoryData models.Categories

	if err = json.NewDecoder(category.Body).Decode(&categoryData); err != nil {
		log.Printf("Failed to decode JSON categories: %v", err)
		categoryChannel <- models.Categories{}
		errChannel <- err
		// close(categoryChannel)
		// close(errChannel)
		return
	}
	categoryChannel <- categoryData
	errChannel <- nil
	// close(categoryChannel)
	// close(errChannel)
	fmt.Println("Category Fetched.")
}

// Fetch all categories from the API.
func FetchCategories(categoriesChannel chan []models.Categories, errChannel chan error) {
	fmt.Println("Starting to Fetch Categories...")
	categories, err := http.Get("http://localhost:3000/api/categories")
	if err != nil {
		log.Printf("Error getting car categories information from the cars API: %v", err)
		categoriesChannel <- nil
		errChannel <- err
		close(categoriesChannel)
		close(errChannel)
		return
	}
	defer categories.Body.Close()

	var categoriesData []models.Categories

	if err = json.NewDecoder(categories.Body).Decode(&categoriesData); err != nil {
		log.Printf("Failed to decode JSON categories: %v", err)
		categoriesChannel <- nil
		errChannel <- err
		close(categoriesChannel)
		close(errChannel)
		return
	}
	categoriesChannel <- categoriesData
	errChannel <- nil
	close(categoriesChannel)
	close(errChannel)
	fmt.Println("Categories fetched.")
}

// Fetch a manufacturer from the API by ID.
func FetchManufacturer(id int, manufacturerChannel chan models.Manufacturers, errChannel chan error) {
	fmt.Printf("Starting to Fetch Manufacturer %d...", id)
	idString := strconv.Itoa(id)

	manufacturer, err := http.Get("http://localhost:3000/api/manufacturers/" + idString)
	if err != nil {
		log.Printf("Error getting car categories information from the cars API: %v", err)
		manufacturerChannel <- models.Manufacturers{}
		errChannel <- err
		// close(manufacturerChannel)
		// close(errChannel)
		return
	}
	defer manufacturer.Body.Close()

	var manufacturerData models.Manufacturers

	if err = json.NewDecoder(manufacturer.Body).Decode(&manufacturerData); err != nil {
		log.Printf("Failed to decode JSON categories: %v", err)
		manufacturerChannel <- models.Manufacturers{}
		errChannel <- err
		// close(manufacturerChannel)
		// close(errChannel)
		return
	}
	manufacturerChannel <- manufacturerData
	errChannel <- nil
	// close(manufacturerChannel)
	// close(errChannel)
	fmt.Println("Manufacturer Fetched.")
}

// Fetch all manufacturers from the API.
func FetchManufacturers(manufacturersChannel chan []models.Manufacturers, errChannel chan error) {
	fmt.Println("Starting to Fetch Manufacturers...")
	manufacturers, err := http.Get("http://localhost:3000/api/manufacturers")
	if err != nil {
		log.Printf("Error getting car categories information from the cars API:%v", err)
		manufacturersChannel <- nil
		errChannel <- err
		close(manufacturersChannel)
		close(errChannel)
		return
	}
	defer manufacturers.Body.Close()

	var manufacturersData []models.Manufacturers

	if err = json.NewDecoder(manufacturers.Body).Decode(&manufacturersData); err != nil {
		log.Printf("Failed to decode JSON categories: %v", err)
		manufacturersChannel <- nil
		errChannel <- err
		close(manufacturersChannel)
		close(errChannel)
		return
	}

	manufacturersChannel <- manufacturersData
	errChannel <- nil
	close(manufacturersChannel)
	close(errChannel)
	fmt.Println("Manufacturers Fetched")
}

// Fetch all models from the API.
func FetchModels(modelsChannel chan []models.Modelcar, errChannel chan error) {
	fmt.Println("Starting to Fetch Models...")
	var carModels []models.Modelcar
	var mu sync.Mutex

	carsDataChannel := make(chan []models.Car)
	carsErrChannel := make(chan error)

	go FetchCars(carsDataChannel, carsErrChannel)
	cars := <-carsDataChannel
	err := <-carsErrChannel

	if err != nil {
		log.Printf("Error Fetching Cars: %v", err)
		modelsChannel <- nil
		errChannel <- err
		return
	}

	for _, car := range cars {
		model := models.Modelcar{
			Id:   car.Id,
			Name: car.Name,
		}
		mu.Lock()
		carModels = append(carModels, model)
		mu.Unlock()
	}
	modelsChannel <- carModels
	errChannel <- nil
	close(modelsChannel)
	close(errChannel)
	fmt.Println("Models Fetched.")
}

// Fetch all. Manufacturers, Categories and Models.
func FetchManCatMod() ([]models.Manufacturers, []models.Categories, []models.Modelcar, error) {
	categoriesChannel := make(chan []models.Categories)
	errCategoriesChannel := make(chan error)
	manufacturersChannel := make(chan []models.Manufacturers)
	errManufacturersChannel := make(chan error)
	modelsChannel := make(chan []models.Modelcar)
	errModelsChannel := make(chan error)

	go FetchCategories(categoriesChannel, errCategoriesChannel)
	go FetchManufacturers(manufacturersChannel, errManufacturersChannel)
	go FetchModels(modelsChannel, errModelsChannel)

	categories := <-categoriesChannel
	err := <-errCategoriesChannel
	if err != nil {
		log.Printf("Error fetching categories: %v", err)
		return nil, nil, nil, err
	}

	manufacturers := <-manufacturersChannel
	err = <-errManufacturersChannel
	if err != nil {
		log.Printf("Error fetching manufacturers: %v", err)
		return nil, nil, nil, err
	}

	dataModels := <-modelsChannel
	err = <-errModelsChannel
	if err != nil {
		log.Printf("Error fetching models: %v", err)
		return nil, nil, nil, err
	}
	return manufacturers, categories, dataModels, nil
}

// Fetch only the cars from the API, that follows the CategoriesFilterMap, ManufacturersFilterMap and ModelsFIlterMap variables.
func FetchFilteredCars() ([]models.Car, error) {
	var carsFiltered []models.Car
	//	Loop through all the cars from the API and collect those whose CategoryID, ManufacturerId and Name have a value of true in
	//	the corresponding maps: CategoriesFilterMap, ManufacturerFilterMap and ModelsFilterMap.
	fmt.Println("Starting FetchFilteredCars....")

	carsDataChannel := make(chan []models.Car, 1)
	errChannel := make(chan error, 1)

	FetchCars(carsDataChannel, errChannel)

	carsData := <-carsDataChannel
	err := <-errChannel
	if err != nil {
		log.Printf("Error Fetching Cars: %v", err)
		return nil, err
	}

	for _, car := range carsData {
		if config.CategoriesFilterMap[car.CategoryID] && config.ManufacturersFilterMap[car.ManufacturerID] && config.ModelsFilterMap[car.Name] {
			carsFiltered = append(carsFiltered, car)
		}
	}
	fmt.Println("FetchFilteredCars DONE")
	return carsFiltered, nil
}

// Fetch only the cars from the API, that are indicated in FavouritesMap variable.
func FetchFavouriteCars() ([]models.Car, error) {
	fmt.Println("Starting to Fetch Favourite Cars...")

	var favouriteCars []models.Car
	var carsSelected []int
	//	Get every ID car from Favourites Map.
	for id, value := range config.FavouritesMap {
		if !value {
			continue
		} else {
			carsSelected = append(carsSelected, id)
		}
	}

	var err error
	var car models.Car
	carsDataChannel := make(chan models.Car, len(carsSelected))
	errChannel := make(chan error, len(carsSelected))

	//	Range trough every ID Car and Fetch the Car.
	for _, carId := range carsSelected {
		go FetchCar(carId, carsDataChannel, errChannel)
	}

	for i := 0; i < len(carsSelected); i++ {
		err = <-errChannel
		car = <-carsDataChannel
		if err != nil {
			return []models.Car{}, err
		} else {
			favouriteCars = append(favouriteCars, car)
		}
	}
	fmt.Println("Favourite Cars Fetched")
	return favouriteCars, nil
}

// Modify the global variable FavouritesMap.
func ModifyFavouritesMap(carId int) {
	config.FavouritesMap[carId] = !config.FavouritesMap[carId]
	fmt.Println("Favourite Map Modified")
}

// Fetch only the cars from the API, that are indicated in ComparisonMap variable.
func FetchComparedCars(compareMap map[int]bool) ([]models.Car, error) {
	var comparedCars []models.Car
	var carsSelected []int

	for id, value := range compareMap {
		if !value {
			continue
		} else {
			carsSelected = append(carsSelected, id)
		}
	}

	var err error
	var car models.Car
	carsDataChannel := make(chan models.Car, len(carsSelected))
	errChannel := make(chan error, len(carsSelected))

	for _, carId := range carsSelected {
		go FetchCar(carId, carsDataChannel, errChannel)
	}

	for i := 0; i < len(carsSelected); i++ {
		err = <-errChannel
		car = <-carsDataChannel
		if err != nil {
			return []models.Car{}, err
		} else {
			comparedCars = append(comparedCars, car)
		}
	}
	fmt.Println("Compared Cars Fetched")
	return comparedCars, nil
}

// Modify the global variable ComparisonMap.
func ModifyComparisonMap(carId int) {

	config.ComparisonMap[carId] = !config.ComparisonMap[carId]
	count := 0
	for _, value := range config.ComparisonMap {
		if value {
			count++
		}
	}
	if count > 1 {
		config.CompareActive = true
		count = 0
	} else {
		config.CompareActive = false
	}
	fmt.Println("Comparison Map Modified")
}

// Resets the global variable ComparisonMap.
func ClearComparisonMap() {

	//	Make all carsID false in the ComparisonFilterMap.
	for key, _ := range config.ComparisonMap {
		config.ComparisonMap[key] = false
	}
	config.CompareActive = false
	fmt.Println("Comparison Map Cleared")
}

// Modify the global variable CategoriesFilterMap.
func ModifyCategoriesFilterMap(selectedItems []string) error {

	//	For those categories selected, make true their values in the CategoriesFilterMap.
	//	Zero Categories selected equals to all categories true.
	if len(selectedItems) > 0 {
		for _, category := range selectedItems {
			categoryId, err := strconv.Atoi(category)
			if err != nil {
				log.Printf("Error. Could not modify the CategoriesFilterMap: %v", err)
				return err
			}
			config.CategoriesFilterMap[categoryId] = true
		}
	} else if len(selectedItems) == 0 {
		for key, _ := range config.CategoriesFilterMap {
			config.CategoriesFilterMap[key] = true
		}
	}
	return nil
}

// Resets the global variable CategoriesFilterMap.
func ClearCategoriesFilterMap() {

	//	Make all categories false in the CategoriesFilterMap
	for key, _ := range config.CategoriesFilterMap {
		config.CategoriesFilterMap[key] = false
	}
}

// Modify the global variable ManufacturersFilterMap.
func ModifyManufacturersFilterMap(selectedItems []string) error {

	//	For those manufacturers selected, make true their values in the ManufacturersFilterMap.
	//	Zero Manufacturers selected equals to all manufacturers true.
	if len(selectedItems) > 0 {
		for _, manufacturer := range selectedItems {
			manufacturerId, err := strconv.Atoi(manufacturer)
			if err != nil {
				log.Printf("Error. Could not modify the ManufacturersFilterMap: %v", err)
				return err
			}
			config.ManufacturersFilterMap[manufacturerId] = true
		}
	} else if len(selectedItems) == 0 {
		for key, _ := range config.ManufacturersFilterMap {
			config.ManufacturersFilterMap[key] = true
		}
	}
	return nil
}

// Resets the global variable ManufacturersFilterMap
func ClearManufacturersFilterMap() {

	//	Make all Manufacturers false in the ManufacturersFilterMap.
	for key, _ := range config.ManufacturersFilterMap {
		config.ManufacturersFilterMap[key] = false
	}
}

// Modify the global variable ModelsFilterMap.
func ModifyModelsFilterMap(selectedItems []string) {

	//	For those models selected, make true their values in the ModelsFilterMap.
	//	Zero Models selected equals to all models true.
	if len(selectedItems) > 0 {
		for _, modelActive := range selectedItems {
			config.ModelsFilterMap[modelActive] = true
		}
	} else if len(selectedItems) == 0 {
		for key, _ := range config.ModelsFilterMap {
			config.ModelsFilterMap[key] = true
		}
	}
}

// Resets the global variable ModelsFilterMap.
func ClearModelsFilterMap() {

	//	Make all Models false in the ModelsFilterMap.
	for key, _ := range config.ModelsFilterMap {
		config.ModelsFilterMap[key] = false
	}
}

func ModifyAllFilterMaps(selectedManufacturers, selectedCategories, selectedModels []string) {
	ModifyManufacturersFilterMap(selectedManufacturers)
	ModifyCategoriesFilterMap(selectedCategories)
	ModifyModelsFilterMap(selectedModels)
}

// Resets all filters: CategoriesFilterMap, ManufacturersFIlterMap and ModelsFilterMap.
func ClearAllFilters() {
	ClearManufacturersFilterMap()
	ClearCategoriesFilterMap()
	ClearModelsFilterMap()
}

// Generates a map with the last comparison data made by the user.
func CreateLastCompareMap() {
	config.LastCompare = make(map[int]bool)
	for key, value := range config.ComparisonMap {
		config.LastCompare[key] = value
	}
}

// func FetchFavouriteCars() ([]models.Car, error) {
// 	fmt.Println("Starting to Fetch Favourite Cars...")
// 	var favouriteCars []models.Car
// 	var wg sync.WaitGroup
// 	var err error
// 	var car models.Car

// 	// Loop through the FavouritesMap and fetch the cars whose ID have a value true on the map.
// 	for id, value := range config.FavouritesMap {
// 		if value {
// 			wg.Add(1)
// 			fmt.Println("Starting go routine.")
// 			carChannel := make(chan models.Car)
// 			carErrChannel := make(chan error)

// 			go func(id int) {
// 				FetchCar(id, carChannel, carErrChannel)
// 				car = <-carChannel
// 				err = <-carErrChannel
// 				defer close(carChannel)
// 				defer close(carErrChannel)
// 			}(id)
// 		} else {
// 			continue
// 		}

// 		if err != nil {
// 			log.Printf("Error Fetching Car %v: %v", id, err)
// 			return []models.Car{}, err
// 		} else {
// 			favouriteCars = append(favouriteCars, car)
// 		}

// 	}

// 	fmt.Println("Favourite Cars Fetched")

// 	return favouriteCars, nil
// }
