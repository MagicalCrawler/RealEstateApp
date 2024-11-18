package utils

import (
	"bytes"
	"encoding/csv"
	"strconv"

	"github.com/MagicalCrawler/RealEstateApp/models"
)

func ExportCSV(input []models.PostHistory) ([]byte, error) {
	buf := new(bytes.Buffer)
	w := csv.NewWriter(buf)
	csvResult := [][]string{
		{
			"title", "url", "price", "deposit", "rent", "city", "neighbor", "area", "bedroom",
			"mode", "type", "age", "floor", "storage", "parking", "elevator",
			"img", "Description", "Capacity", "NormalDays", "weekend", "holidays", "CostPerPerson",
		},
	}
	for _, post := range input {
		csvResult = append(csvResult, []string{
			post.Title, post.PostURL, strconv.FormatInt(post.Price, 10), strconv.FormatInt(post.Deposit, 10),
			strconv.FormatInt(post.Rent, 10), post.City, post.Neighborhood, strconv.Itoa(post.Area), strconv.Itoa(post.BedroomNum),
			string(post.BuyMode), string(post.Building), strconv.Itoa(int(post.Age)), strconv.Itoa(int(post.FloorsNum)),
			strconv.FormatBool(post.HasStorage), strconv.FormatBool(post.HasParking), strconv.FormatBool(post.HasElevator),
			post.ImageURL, post.Description, post.Capacity, post.NormalDays, post.Weekend, post.Holidays, post.CostPerPerson,
		})
	}
	w.WriteAll(csvResult)
	if err := w.Error(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}