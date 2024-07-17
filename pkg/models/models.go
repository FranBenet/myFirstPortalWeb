package models

type Car struct {
	Id             int    `json:"id"`
	Name           string `json:"name"`
	ManufacturerID int    `json:"manufacturerId"`
	CategoryID     int    `json:"categoryId"`
	Year           int    `json:"year"`
	Specifications Specs  `json:"specifications"`
	Image          string `json:"image"`
}

type Manufacturers struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	Country      string `json:"country"`
	FoundingYear int    `json:"foundingYear"`
}

type Categories struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Modelcar struct {
	Id   int
	Name string
}

type Specs struct {
	Engine       string `json:"engine"`
	Horsepower   int    `json:"horsepower"`
	Transmission string `json:"transmission"`
	DriveTrain   string `json:"drivetrain"`
}

// Card is the struct created for the Gallery on the main page
type Card struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	Manufacturer string
	Category     string
	Year         int    `json:"year"`
	Image        string `json:"image"`
	Liked        bool
	Compared     bool
}

// ExtendedCard is the struct created for when a car is clicked, or when viewing the favourites or compare pages.
type ExtendedCard struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	Manufacturer string
	Country      string `json:"country"`
	FoundingYear int    `json:"foundingYear"`
	Category     string
	Year         int    `json:"year"`
	Engine       string `json:"engine"`
	Horsepower   int    `json:"horsePower"`
	Transmission string `json:"transmission"`
	DriveTrain   string `json:"driveTrain"`
	Image        string `json:"image"`
	Liked        bool
	Compared     bool
}

// DataResponse is the struct used to send in the response with the HTML.
type DataResponse struct {
	Card          []Card
	ExtCard       []ExtendedCard
	Manufacturers []Manufacturers
	Categories    []Categories
	Models        []Modelcar
	NoResults     bool
	CompareActive bool
}

type CarSearch struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Manufacture string `json:"manufacturerId"`
	Category    string `json:"categoryId"`
}
