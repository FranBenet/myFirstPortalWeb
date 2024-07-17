package config

var FavouritesMap map[int]bool
var ComparisonMap map[int]bool
var RedirectURL string
var CompareActive bool

var ManufacturersFilterMap map[int]bool
var CategoriesFilterMap map[int]bool
var ModelsFilterMap map[string]bool
var TotalNumCars = 0
var LastCompare map[int]bool

func init() {
	FavouritesMap = make(map[int]bool)
	ComparisonMap = make(map[int]bool)
	RedirectURL = "/"
	CompareActive = false

	ManufacturersFilterMap = make(map[int]bool)
	CategoriesFilterMap = make(map[int]bool)
	ModelsFilterMap = make(map[string]bool)
}
