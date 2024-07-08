package models

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

type CarModels struct {
	Id             int    `json:"id"`
	Name           string `json:"name"`
	ManufacturerID int    `json:"manufacturerId"`
	CategoryID     int    `json:"categoryId"`
	Year           int    `json:"year"`
	Specs          Specs  `json:"specs"`
	Image          string `json:"image"`
}

type Specs struct {
	Engine       string `json:"engine"`
	Horsepower   int    `json:"horsePower"`
	Transmission string `json:"transmission"`
	DriveTrain   string `json:"driveTrain"`
}
