package exporter

import (
	"strings"
	"testing"

	"github.com/MagicalCrawler/RealEstateApp/models"
	"github.com/MagicalCrawler/RealEstateApp/types"
	"github.com/MagicalCrawler/RealEstateApp/utils"
)

func TestExportCSV(t *testing.T) {
	p1 := models.PostHistory{ID: 1, PostID: 1, Title: "p1", PostURL: "url1", Price: 11, Deposit: 12, Rent: 13, City: "Tehran",
		Neighborhood: "n1", Area: 14, BedroomNum: 15, BuyMode: types.Rent, Building: types.Apartment, Age: 15, FloorsNum: 16, HasStorage: true, HasParking: false,
		HasElevator: true, ImageURL: "i-url1", Description: "d1", Capacity: "c1", NormalDays: "nd1", Weekend: "w1", Holidays: "h1", CostPerPerson: "cpp1",
	}
	p2 := models.PostHistory{ID: 2, PostID: 2, Title: "p2", PostURL: "url2", Price: 21, Deposit: 22, Rent: 23, City: "Tehran",
		Neighborhood: "n2", Area: 24, BedroomNum: 25, BuyMode: types.Rent, Building: types.Apartment, Age: 25, FloorsNum: 26, HasStorage: true, HasParking: false,
		HasElevator: true, ImageURL: "i-url2", Description: "d2", Capacity: "c2", NormalDays: "nd2", Weekend: "w2", Holidays: "h2", CostPerPerson: "cpp2",
	}
	postHistories := []models.PostHistory{p1, p2}
	bytesResult, err := utils.ExportCSV(postHistories)
	if err != nil {
		t.Errorf("export csv failed: %v", err)
	}
	result := string(bytesResult)
	resultLines := strings.Split(strings.TrimSpace(string(result)), "\n")
	expectedLines := []string{
		"title,url,price,deposit,rent,city,neighbor,area,bedroom,mode,type,age,floor,storage,parking,elevator,img,Description,Capacity,NormalDays,weekend,holidays,CostPerPerson",
		"p1,url1,11,12,13,Tehran,n1,14,15,rent,apartment,15,16,true,false,true,i-url1,d1,c1,nd1,w1,h1,cpp1",
		"p2,url2,21,22,23,Tehran,n2,24,25,rent,apartment,25,26,true,false,true,i-url2,d2,c2,nd2,w2,h2,cpp2",
	}
	if len(resultLines) != len(expectedLines) {
		t.Error("result length is not matched")
	}
	for index, expectedLine := range expectedLines {
		if resultLines[index] != expectedLine {
			t.Error("result is not matched:")
			t.Logf("result: %v", resultLines[index])
			t.Logf("expected: %v", expectedLine)
		}
	}
}
