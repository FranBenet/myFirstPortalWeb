package handlers

import (
	"cars/pkg/config"
	"cars/pkg/helpers"
	"cars/pkg/models"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

// Responds with the index page including the gallery of all the cars from the API.
func Homepage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		log.Printf("Error. Homepage - Path Not Allowed")
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
		log.Println("Error. Homepage - Method Not Allowed")
		return
	}
	fmt.Println("Starting Homepage...")

	//	We store the current URL
	config.RedirectURL = r.URL.String()

	carsDataChannel := make(chan []models.Car, 1)
	errCarsChannel := make(chan error, 1)

	go helpers.FetchCars(carsDataChannel, errCarsChannel)

	carsData := <-carsDataChannel
	err := <-errCarsChannel
	if err != nil {
		// log.Printf("Error fetching cars: %v", err)
		// http.Error(w, "Failed to fetch cars data", http.StatusInternalServerError)
		// return
	}

	//	Create Small Card for each car.
	cards, err := helpers.CreateSmallCardsBatch(carsData)
	if err != nil {
		// log.Println("Error Creating Big Cars.")
	}
	fmt.Println("Cards Created.")

	manufacturers, categories, dataModels, err := helpers.FetchManCatMod()
	if err != nil {
		log.Printf("Error Fetching Data from the API:%v", err)
	}

	//	Create the variable to store the data that needs to be send in our HTML Response.
	//	Add the cards information, the list of categories and manufacturers to their specific key.
	var data models.DataResponse
	data.Card = cards
	data.Categories = categories
	data.Manufacturers = manufacturers
	data.Models = dataModels
	data.NoResults = false
	data.CompareActive = config.CompareActive

	htmlTemplates := []string{
		"web/templates/index.html",
		"web/templates/main-bar.html",
		"web/templates/filter.html",
		"web/templates/card-template.html",
	}

	tmpl, err := template.ParseFiles(htmlTemplates...)
	if err != nil {
		log.Printf("Error. Homepage - Parsing Template: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		log.Printf("Error. Homepage - Executing Template: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	fmt.Println("HOMEPAGE DONE")
	fmt.Println()
	fmt.Println()
	fmt.Println()
	fmt.Println()
	fmt.Println("-------------------------------")
}

// Responds with a page including only the car selected by the user. This includes extra information.
func SelectCar(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/id" {
		log.Printf("Error. SelectCar - Path Not Allowed")
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		log.Println("Error. SelectCar - Method Not Allowed")
		http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
		return
	}

	//	Getting QUERY key=map from the URL
	//	(HTML allows to send this type of query: ?id=id_number to be accessed in the server.
	//	The query is after the ?. The = symbol divides the key from the value.)
	carID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		log.Printf("Error selecting car. Could not convert query into Integer: %v", err)
		return
	}
	// We store the current URL
	config.RedirectURL = r.URL.String()

	//	Fetch the selected car.
	carDataChannel := make(chan models.Car, 1)
	errChannel := make(chan error, 1)

	helpers.FetchCar(carID, carDataChannel, errChannel)

	carData := <-carDataChannel
	err = <-errChannel
	close(carDataChannel)
	close(errChannel)
	if err != nil {
		log.Printf("Error fetching car %v: %v", carID, err)
		return
	}

	if carID > config.TotalNumCars {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	//	Create a Big Card version of the Car that was selected.
	card, err := helpers.CreateBigCard(carData)
	if err != nil {
		log.Printf("Error Creating Big Card: %v", err)
		http.Error(w, "Failed to create big cards", http.StatusInternalServerError)
		return
	}

	//	Create the variable to be sent with the HTML.
	//	And add the data on it.
	var data models.DataResponse
	data.ExtCard = append(data.ExtCard, card)
	data.CompareActive = config.CompareActive

	htmlTemplates := []string{
		"web/templates/card-page.html",
		"web/templates/main-bar.html",
		"web/templates/card-template.html",
	}

	tmpl, err := template.ParseFiles(htmlTemplates...)
	if err != nil {
		log.Printf("Error. SelectCar - Parsing Template: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "card-page.html", data)
	if err != nil {
		log.Printf("Error. SelectCar - Executing Template: %v\n", err)
		http.Error(w, "Internal Server Error ID", http.StatusInternalServerError)
		return
	}

}

// Modifies the FavouritesMap and CompareMap and Responds by redirecting to the URL where it came from.
func StatusChange(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/liked-compared" {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		log.Printf("Error. StatusChange - Path Not Allowed")
		return
	}
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
		log.Printf("Error. StatusChange - Method Not Allowed")
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		log.Printf("Error. StatusChange - Parsing Form")
		return
	}

	//	Variables that extract the information from the form that is included in each Card.
	//	form_id -> the ID of the car that triggered the request.
	//	trigger -> the information about what button was selected: Favourite or Compare
	//	redirect_url -> information about the current URL for the user, when it submitted the form.
	carId, err := strconv.Atoi(r.Form.Get("form_id"))
	if err != nil {
		log.Printf("Error Status Change. Could not convert query into Integer: %v", err)
		return
	}
	triggeredButton := r.Form.Get("trigger")

	//	Check what button from the form was selected and change the corresponding map.
	if triggeredButton == "favorite" {
		helpers.ModifyFavouritesMap(carId)
	} else if triggeredButton == "compare" {
		helpers.ModifyComparisonMap(carId)
	} else {
		http.Redirect(w, r, config.RedirectURL, http.StatusInternalServerError)
	}

	//	Redirect the client to the URL where the Request came from.
	http.Redirect(w, r, config.RedirectURL, http.StatusSeeOther)
}

// Responds with the compare page including the cars selected to be compared.
func ComparePage(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/comparePage" {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		log.Printf("Error. ComparisonPage - Path Not Allowed")
		return
	}
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
		log.Printf("Error. ComparisonPage - Method Not Allowed")
		return
	}

	currentURL := r.URL.String()
	var comparedCars []models.Car
	var err error

	//	Collect the cars from different maps, depending on the URL where the request comes from.
	// This is done, to allow user "like" a car from the comparison page.
	if config.RedirectURL != currentURL {

		comparedCars, err = helpers.FetchComparedCars(config.ComparisonMap)
		if err != nil {
			log.Printf("Error fetching compared cars: %v", err)
			return
		}
		helpers.CreateLastCompareMap()
		helpers.ClearComparisonMap()
		config.RedirectURL = r.URL.String()

	} else if config.RedirectURL == currentURL {

		comparedCars, err = helpers.FetchComparedCars(config.LastCompare)
		if err != nil {
			log.Printf("Error fetching compared cars: %v", err)
			return
		}
	}

	if len(comparedCars) < 2 {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		//	Create Big Card for each car.
		cards, err := helpers.CreateBigCardsBatch(comparedCars)
		if err != nil {
			log.Println("Error Creating Big Cars.")
		}

		//	Create a variable to be sent together with the HTML.
		//	Add the data from the car/s on it
		var data models.DataResponse
		data.ExtCard = cards
		data.CompareActive = config.CompareActive

		htmlTemplates := []string{
			"web/templates/card-page.html",
			"web/templates/main-bar.html",
			"web/templates/card-template.html",
		}

		tmpl, err := template.ParseFiles(htmlTemplates...)
		if err != nil {
			log.Printf("Error. ComparisonPage - Parsing Template: %v\n", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		err = tmpl.ExecuteTemplate(w, "card-page.html", data)
		if err != nil {
			log.Printf("Error. ComparisonPage - Executing Template: %v\n", err)
			http.Error(w, "Internal Server Error ID", http.StatusInternalServerError)
			return
		}
	}
}

// Responds with a page including all the cars that have been liked.
func FavouritesPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/favouritePage" {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		log.Printf("Error. FavouritesPage - Path Not Allowed")
		return
	}
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
		log.Printf("Error. FavouritesPage - Method Not Allowed")
		return
	}

	// We store the current URL
	config.RedirectURL = r.URL.String()

	favouriteCars, err := helpers.FetchFavouriteCars()
	if err != nil {
		log.Printf("Error fetching favourite cars: %v", err)
		return
	}

	if len(favouriteCars) == 0 {
		NoResultsCardPage(w)
	} else {
		//	Create Big Card for each car.
		cards, err := helpers.CreateBigCardsBatch(favouriteCars)
		if err != nil {
			log.Println("Error Creating Big Cars.")
		}

		//	Create a variable to be sent together with the HTML.
		//	Add the data from the car/s on it
		var data models.DataResponse
		data.ExtCard = cards
		data.CompareActive = config.CompareActive

		htmlTemplates := []string{
			"web/templates/card-page.html",
			"web/templates/main-bar.html",
			"web/templates/card-template.html",
		}

		tmpl, err := template.ParseFiles(htmlTemplates...)
		if err != nil {
			log.Printf("Error. FavouritesPage - Parsing Template: %v\n", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		err = tmpl.ExecuteTemplate(w, "card-page.html", data)
		if err != nil {
			log.Printf("Error. FavouritesPage - Executing Template: %v\n", err)
			http.Error(w, "Internal Server Error ID", http.StatusInternalServerError)
			return
		}
	}
}

// Responds with the index page but without any cars. A message "0 results found" instead will be shown.
func NoResultsIndex(w http.ResponseWriter) {

	manufacturers, categories, dataModels, err := helpers.FetchManCatMod()

	var data models.DataResponse
	data.Categories = categories
	data.Manufacturers = manufacturers
	data.Models = dataModels
	data.NoResults = true
	data.CompareActive = config.CompareActive

	htmlTemplates := []string{
		"web/templates/index.html",
		"web/templates/main-bar.html",
		"web/templates/filter.html",
		"web/templates/card-template.html",
	}

	tmpl, err := template.ParseFiles(htmlTemplates...)
	if err != nil {
		log.Printf("Error. NoResultsIndex - Parsing Template: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		log.Printf("Error. NoResultsIndex - Executing Template: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// Responds with the card-page but without any cars. A message "0 results found" instead will be shown.
func NoResultsCardPage(w http.ResponseWriter) {
	var data models.DataResponse
	data.NoResults = true
	data.CompareActive = config.CompareActive

	htmlTemplates := []string{
		"web/templates/card-page.html",
		"web/templates/main-bar.html",
		"web/templates/card-template.html",
	}
	tmpl, err := template.ParseFiles(htmlTemplates...)
	if err != nil {
		log.Printf("Error. NoResultsCardPage - Parsing Template: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "card-page.html", data)
	if err != nil {
		log.Printf("Error. NoResultsCardPage - Executing Template: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

}

// Responds with the last compare made by the user.
func LastCompare(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/lastCompare" {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		log.Printf("Error. ComparisonPage - Path Not Allowed")
		return
	}
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
		log.Printf("Error. ComparisonPage - Method Not Allowed")
		return
	}

	// We store the current URL
	config.RedirectURL = r.URL.String()

	//	FetchComparedCars collects all the cars marked to be Compared
	//	and stores them in comparedCars variable.
	comparedCars, err := helpers.FetchComparedCars(config.LastCompare)

	if err != nil {
		log.Printf("Error finding last compared cars: %v", err)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if len(comparedCars) == 0 {
		NoResultsCardPage(w)
	} else {
		//	Create Big Card for each car.
		cards, err := helpers.CreateBigCardsBatch(comparedCars)
		if err != nil {
			log.Println("Error Creating Big Cars.")
		}

		//	Create a variable to be sent together with the HTML.
		//	Add the data from the car/s on it
		var data models.DataResponse
		data.ExtCard = cards
		data.CompareActive = config.CompareActive

		htmlTemplates := []string{
			"web/templates/card-page.html",
			"web/templates/main-bar.html",
			"web/templates/card-template.html",
		}

		tmpl, err := template.ParseFiles(htmlTemplates...)
		if err != nil {
			log.Printf("Error. Last Compare Page - Parsing Template: %v\n", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		err = tmpl.ExecuteTemplate(w, "card-page.html", data)
		if err != nil {
			log.Printf("Error. Last Comapre Page - Executing Template: %v\n", err)
			http.Error(w, "Internal Server Error ID", http.StatusInternalServerError)
			return
		}
	}
}

// Takes the filters or search request and Responds with the cars that matches.
func Filter(w http.ResponseWriter, r *http.Request) {
	fmt.Println("FILTER PAGE STARTING...")
	fmt.Println()

	if r.URL.Path != "/search" {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		log.Printf("Error. Filter - Path Not Allowed")
		return
	}
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
		log.Printf("Error. Filter - Method Not Allowed")
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		log.Printf("Error. Filter - Parsing Form")
		return
	}

	helpers.ClearAllFilters()
	fmt.Println("All filter cleaned.")

	config.RedirectURL = r.URL.String()
	//	Get the Form input. This allows us to know what button triggered the submission.
	action := r.FormValue("action")

	//	In case there is no string in the action. That means the search bar triggered the submission.
	if action == "" {
		//	We get the query from the searchbar and check if it has content or not.
		query := r.FormValue("searchRequest")
		if query != "" {
			filteredCars, err := helpers.SearchQueryCars(query)
			if err != nil {
				log.Printf("Error searching cars from the query", err)
				//	WE DO SOMETHING
			}
			//	Check the number of cars fetched to determine whether we display a
			//	"0 results found" or not.
			if len(filteredCars) == 0 {
				NoResultsIndex(w)
			} else {
				cards, err := helpers.CreateSmallCardsBatch(filteredCars)
				if err != nil {
					log.Printf("Error Creating Cards:", err)
				}

				manufacturers, categories, dataModels, err := helpers.FetchManCatMod()
				if err != nil {
					log.Printf("Error fetching categories: %v", err)
					http.Error(w, "Failed to fetch items from the APO:", http.StatusInternalServerError)
				}

				//	Create a variable to be sent together with the HTML.
				//	Add the data from the car/s on it
				var data models.DataResponse
				data.Card = cards
				data.Categories = categories
				data.Manufacturers = manufacturers
				data.Models = dataModels
				data.CompareActive = config.CompareActive

				htmlTemplates := []string{
					"web/templates/index.html",
					"web/templates/main-bar.html",
					"web/templates/filter.html",
					"web/templates/card-template.html",
				}

				tmpl, err := template.ParseFiles(htmlTemplates...)
				if err != nil {
					log.Printf("Error. Filters - Parsing Template: %v\n", err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}

				err = tmpl.ExecuteTemplate(w, "index.html", data)
				if err != nil {
					log.Printf("Error. Filters - Executing Template: %v\n", err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}
			}
		} else {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
		//	If the action contains some string. It means the submission
		//	was triggered by the filter buttons.
	} else if action != "" {
		fmt.Println("Action is something")
		selectedManufacturers := r.Form["manufacturer"]
		selectedCategories := r.Form["category"]
		selectedModels := r.Form["model"]

		switch action {
		case "search":
			helpers.ModifyAllFilterMaps(selectedManufacturers, selectedCategories, selectedModels)

		case "acceptManufacturer":
			helpers.ModifyAllFilterMaps(selectedManufacturers, selectedCategories, selectedModels)

		case "clearManufacturer":
			helpers.ClearManufacturersFilterMap()

		case "acceptCategory":
			helpers.ModifyAllFilterMaps(selectedManufacturers, selectedCategories, selectedModels)

		case "clearCategory":
			helpers.ClearCategoriesFilterMap()

		case "acceptModel":
			helpers.ModifyAllFilterMaps(selectedManufacturers, selectedCategories, selectedModels)

		case "clearModel":
			helpers.ClearModelsFilterMap()
		}
		fmt.Println("Filter Maps Modified")
		//	We fetch the filtered Cars
		filteredCars, err := helpers.FetchFilteredCars()
		if err != nil {
			log.Printf("Error fetching filtered cars: %v", err)
			return
		}

		//If no filteredCars -> Print: No results page
		if len(filteredCars) == 0 {
			NoResultsIndex(w)
		} else {
			//	Create for each car a small card.
			cards, err := helpers.CreateSmallCardsBatch(filteredCars)
			if err != nil {
				log.Println("Error Creating Smmall Cards")
			}

			//	Fetch manufacturers, categories and models.
			manufacturers, categories, dataModels, err := helpers.FetchManCatMod()
			if err != nil {
				log.Printf("Error fetching categories: %v", err)
				http.Error(w, "Failed to fetch items from the APO:", http.StatusInternalServerError)
			}

			//	Create a variable to be sent together with the HTML.
			//	Add the data from the car/s on it
			var data models.DataResponse
			data.Card = cards
			data.Categories = categories
			data.Manufacturers = manufacturers
			data.Models = dataModels
			data.CompareActive = config.CompareActive

			htmlTemplates := []string{
				"web/templates/index.html",
				"web/templates/main-bar.html",
				"web/templates/filter.html",
				"web/templates/card-template.html",
			}

			tmpl, err := template.ParseFiles(htmlTemplates...)
			if err != nil {
				log.Printf("Error. Filters - Parsing Template: %v\n", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			err = tmpl.ExecuteTemplate(w, "index.html", data)
			if err != nil {
				log.Printf("Error. Filters - Executing Template: %v\n", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}
	}
	fmt.Println("FILTER DONE")
	fmt.Println("-------------------------------")
	fmt.Println()
	fmt.Println()
	fmt.Println()
	fmt.Println()
}
