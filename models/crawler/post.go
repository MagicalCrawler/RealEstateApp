package crawlerModels

import "github.com/MagicalCrawler/RealEstateApp/types"

// RentalMetadata holds rental-specific details for posts
type RentalMetadata struct {
	Capacity        string
	NormalDayPrice  string
	WeekendPrice    string
	HolidayPrice    string
	ExtraPersonCost string
}

// Post represents a real estate listing
type Post struct {
	ID                  string
	Title               string
	City                City
	Neighborhood        string
	Price               string
	Link                string
	Images              []string
	Description         string
	Area                string
	YearBuilt           string
	Rooms               string
	PricePerSquareMeter string
	TotalPrice          string
	Floor               string
	Features            []string
	Deposit             string
	MonthlyRent         string
	DepositOnRentDesc   string
	RentalMetadata      *RentalMetadata
	Website             types.WebsiteSource
}

// City represents a city in the system
type City struct {
	Name  string `json:"name"`
	Slug  string `json:"slug"`
	Level string `json:"level"`
}

// CityResponse is used to parse the API response for cities
type CityResponse struct {
	Cities []City `json:"cities"`
}
