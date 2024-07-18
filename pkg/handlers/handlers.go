package handlers

import (
	"cars/pkg/config"
	"cars/pkg/helpers"
	"cars/pkg/models"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

// Responds with the index page including the gallery of all the cars from the API.
func Homepage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		fmt.Println("Error. Path Not Allowed. Homepage")
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
		return
	}

	//	We store the current URL to keep track of redirection when needed.
	config.RedirectURL = r.URL.String()

	//	Fetch the cars from the API
	carsDataChannel := make(chan []models.Car, 1)
	errCarsChannel := make(chan error, 1)

	go helpers.FetchCars(carsDataChannel, errCarsChannel)

	carsData := <-carsDataChannel
	err := <-errCarsChannel
	close(carsDataChannel)
	close(errCarsChannel)

	if err != nil {
		fmt.Println("Error fetching data from the API.")
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, "Error fetching data from the API.", http.StatusInternalServerError)
		return
	}

	//	Create a small card for each car. Small Card just refers to a variable with sjust few data ot the cars.
	cards, err := helpers.CreateSmallCardsBatch(carsData)
	if err != nil {
		fmt.Println("Error Creating  cards.")
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, "Error Creating  cards.", http.StatusInternalServerError)
		return
	}

	//	Fetch data from the API related to manufacturers, categories and models.
	manufacturers, categories, dataModels, err := helpers.FetchManCatMod()
	if err != nil {
		fmt.Println("Error fetching data from the API.")
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, "Error fetching data from the API.", http.StatusInternalServerError)
		return
	}

	//	Collect the data to be send with the HTML
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

	helpers.RenderTemplate(w, htmlTemplates, "index.html", data)

}

// Responds with a page including only the car selected by the user. This includes extra information.
func SelectCar(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/id" {
		fmt.Println("Error. Path Not Allowed. ID")
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
		return
	}

	//	Getting QUERY key=map from the URL
	//	(HTML allows to send this type of query: ?id=id_number to be accessed in the server.
	//	The query is after the ?. The = symbol divides the key from the value.
	carID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		fmt.Println("Error selecting car. Could not convert query into Integer: ", err)
		return
	}

	//	We store the current URL to keep track of redirection when needed.
	config.RedirectURL = r.URL.String()

	//	Fetch the selected car from the API.
	carDataChannel := make(chan models.Car, 1)
	errChannel := make(chan error, 1)

	helpers.FetchCar(carID, carDataChannel, errChannel)

	carData := <-carDataChannel
	err = <-errChannel
	close(carDataChannel)
	close(errChannel)
	if err != nil {
		fmt.Println("Error fetching car from the API.")
		w.WriteHeader(http.StatusInternalServerError)
		NotFoundHandler(w, r)
		return
	}

	//	We handle the situation where a URL is added with a non-existent car ID by redirecting to main.page.
	if carID > config.TotalNumCars {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	//	Create a big card for the selected car. Big cards refers to a variable including more data than the one included in the small cards.
	card, err := helpers.CreateBigCard(carData)
	if err != nil {
		fmt.Println("Error Creating card.")
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, "Error Creating card.", http.StatusInternalServerError)
		return
	}

	//	Create the variable to be sent with the HTML and add the data on it.
	var data models.DataResponse
	data.ExtCard = append(data.ExtCard, card)
	data.CompareActive = config.CompareActive

	htmlTemplates := []string{
		"web/templates/card-page.html",
		"web/templates/main-bar.html",
		"web/templates/card-template.html",
	}

	helpers.RenderTemplate(w, htmlTemplates, "card-page.html", data)

}

// Modifies the FavouritesMap and CompareMap and Responds by redirecting to the URL where it came from.
func StatusChange(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/liked-compared" {

		fmt.Println("Error. Path Not Allowed. StatusChange")
		fmt.Println("Is this the error?")
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		fmt.Println("Error Parsing Form")
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	//	Variables that extract the information from the form that is included in each Card.
	//	form_id -> the ID of the car that triggered the request.
	//	trigger -> the information about what button was selected: Favourite or Compare
	//	redirect_url -> information about the URL where the user clicked the button when submitted the form.
	carId, err := strconv.Atoi(r.Form.Get("form_id"))
	if err != nil {
		fmt.Println("Error converting form_id.")
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, "Error converting form_id.", http.StatusInternalServerError)
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
		return
	}

	//	Redirect the client to the URL where the Request came from.
	http.Redirect(w, r, config.RedirectURL, http.StatusSeeOther)

}

// Responds with the compare page including the cars selected to be compared.
func ComparePage(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/comparePage" {
		fmt.Println("Error. Path Not Allowed. Comparepage")
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
		return
	}

	currentURL := r.URL.String()
	var comparedCars []models.Car
	var err error

	//	Collect the cars from different maps, depending if the request comes from "Last comparison" button or "Compare button".
	// 	Checking the URL is done to allow user "like" a car from the comparison page as well.
	if config.RedirectURL != currentURL {

		comparedCars, err = helpers.FetchComparedCars(config.ComparisonMap)
		if err != nil {
			fmt.Println("Error fetching data from the API.")
			w.WriteHeader(http.StatusInternalServerError)
			NotFoundHandler(w, r)
			return
		}
		helpers.CreateLastCompareMap()
		helpers.ClearComparisonMap()
		config.RedirectURL = r.URL.String()

	} else if config.RedirectURL == currentURL {

		comparedCars, err = helpers.FetchComparedCars(config.LastCompare)
		if err != nil {
			fmt.Println("Error fetching compared cars.")
			w.WriteHeader(http.StatusInternalServerError)
			NotFoundHandler(w, r)
			return
		}
	}

	//	Check that actually there are some cars to be display. Otherwise, redirect to main page.
	if len(comparedCars) < 2 {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		//	Create Big Card for each car.
		cards, err := helpers.CreateBigCardsBatch(comparedCars)
		if err != nil {
			fmt.Println("Error creating cards.")
			w.WriteHeader(http.StatusInternalServerError)
			http.Error(w, "Error creating cards.", http.StatusInternalServerError)
			return
		}

		//	Create a variable to be sent together with the HTML.
		//	Add the data from the cars on it
		var data models.DataResponse
		data.ExtCard = cards
		data.CompareActive = config.CompareActive

		htmlTemplates := []string{
			"web/templates/card-page.html",
			"web/templates/main-bar.html",
			"web/templates/card-template.html",
		}

		helpers.RenderTemplate(w, htmlTemplates, "card-page.html", data)
	}
}

// Responds with a page including all the cars that have been liked.
func FavouritesPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/favouritePage" {
		fmt.Println("Error. Path Not Allowed. FavouritePage")
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
		return
	}

	// We store the current URL
	config.RedirectURL = r.URL.String()

	favouriteCars, err := helpers.FetchFavouriteCars()
	if err != nil {
		fmt.Println("Error fetching data from the API.")
		w.WriteHeader(http.StatusInternalServerError)
		NotFoundHandler(w, r)
		return
	}

	if len(favouriteCars) == 0 {
		NoResultsCardPage(w)
	} else {
		//	Create Big Card for each car.
		cards, err := helpers.CreateBigCardsBatch(favouriteCars)
		if err != nil {
			fmt.Println("Error creating cards.")
			w.WriteHeader(http.StatusInternalServerError)
			http.Error(w, "Error creating cards.", http.StatusInternalServerError)
			return
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

		helpers.RenderTemplate(w, htmlTemplates, "card-page.html", data)
	}
}

// Responds with the index page but without any cars. A message "0 results found" instead will be shown.
func NoResultsIndex(w http.ResponseWriter) {

	manufacturers, categories, dataModels, err := helpers.FetchManCatMod()
	if err != nil {
		fmt.Println("Error fetching data from the API.")
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, "Error fetching data from the API.", http.StatusInternalServerError)
		return
	}

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

	helpers.RenderTemplate(w, htmlTemplates, "index.html", data)
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
	helpers.RenderTemplate(w, htmlTemplates, "card-page.html", data)

}

// Responds with the last compare made by the user.
func LastCompare(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/lastCompare" {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		fmt.Println("Error Path Not Allowed. ComparePage")
		return
	}
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
		return
	}

	// We store the current URL
	config.RedirectURL = r.URL.String()

	//	FetchComparedCars collects all the cars marked to be Compared
	//	and stores them in comparedCars variable.
	comparedCars, err := helpers.FetchComparedCars(config.LastCompare)
	if err != nil {
		fmt.Println("Error finding last compared cars: ", err)
		NotFoundHandler(w, r)
		return
	}

	if len(comparedCars) == 0 {
		NoResultsCardPage(w)
	} else {
		//	Create Big Card for each car.
		cards, err := helpers.CreateBigCardsBatch(comparedCars)
		if err != nil {
			fmt.Println("Error creating cards.")
			w.WriteHeader(http.StatusInternalServerError)
			http.Error(w, "Error creating cards.", http.StatusInternalServerError)
			return
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

		helpers.RenderTemplate(w, htmlTemplates, "card-page.html", data)
	}
}

// Takes the filters or search request and Responds with the cars that matches.
func Filter(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/search" {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		fmt.Println("Error Path Not Allowed. FilterPage")
		return
	}
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		fmt.Println("Error Parsing Form")
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	helpers.ClearAllFilters()

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
				fmt.Println("Error filtering data.")
				w.WriteHeader(http.StatusInternalServerError)
				NotFoundHandler(w, r)
				return
			}
			//	Check the number of cars fetched to determine whether we display a
			//	"0 results found" or not.
			if len(filteredCars) == 0 {
				NoResultsIndex(w)
			} else {
				cards, err := helpers.CreateSmallCardsBatch(filteredCars)
				if err != nil {
					fmt.Println("Error creating cards.")
					w.WriteHeader(http.StatusInternalServerError)
					http.Error(w, "Error creating cards.", http.StatusInternalServerError)
					return
				}

				manufacturers, categories, dataModels, err := helpers.FetchManCatMod()
				if err != nil {
					fmt.Println("Error fetching data from the API.")
					w.WriteHeader(http.StatusInternalServerError)
					NotFoundHandler(w, r)
					return
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

				helpers.RenderTemplate(w, htmlTemplates, "index.html", data)
			}
		} else {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
		//	If the action contains some string. It means the submission
		//	was triggered by the filter buttons.
	} else if action != "" {
		selectedManufacturers := r.Form["manufacturer"]
		selectedCategories := r.Form["category"]
		selectedModels := r.Form["model"]

		helpers.ModifyAllFilterMaps(selectedManufacturers, selectedCategories, selectedModels)

		//	We fetch the filtered Cars
		filteredCars, err := helpers.FetchFilteredCars()
		if err != nil {
			fmt.Println("Error filtering data.")
			w.WriteHeader(http.StatusInternalServerError)
			NotFoundHandler(w, r)
			return
		}

		//If no filteredCars -> Print: No results page
		if len(filteredCars) == 0 {
			NoResultsIndex(w)
		} else {
			//	Create for each car a small card.
			cards, err := helpers.CreateSmallCardsBatch(filteredCars)
			if err != nil {
				fmt.Println("Error creating cards.")
				w.WriteHeader(http.StatusInternalServerError)
				http.Error(w, "Error creating cards.", http.StatusInternalServerError)
				return
			}

			//	Fetch manufacturers, categories and models.
			manufacturers, categories, dataModels, err := helpers.FetchManCatMod()
			if err != nil {
				fmt.Println("Error fetching data from the API.")
				w.WriteHeader(http.StatusInternalServerError)
				NotFoundHandler(w, r)
				return
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

			helpers.RenderTemplate(w, htmlTemplates, "index.html", data)
		}
	}
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {

	htmlTemplates := []string{
		"web/templates/500.html",
	}

	tmpl, err := template.ParseFiles(htmlTemplates...)
	if err != nil {
		fmt.Println("Error Parsing 500HTML Template: ", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, "500.html")
	if err != nil {
		fmt.Println("Error Executing 500HTML Template: ", err)
		http.Error(w, "Internal Server Error ID", http.StatusInternalServerError)
		return
	}

}
