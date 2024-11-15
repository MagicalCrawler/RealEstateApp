package utils

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func GetConfig(key string) string {
	return os.Getenv(key)
}

func LoadEnvFile() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println(err)
		panic("Error loading .env file")
	}
}

type ProvincialCenter struct {
	Name string `json:"name"`
}

type AppSettings struct {
	ProvincialCenters []ProvincialCenter `json:"Provincial-Centers"`
}

func LoadAppSettingsFile() ([]ProvincialCenter, error) {
	file, err := os.Open("appsettings.json")
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var settings AppSettings
	if err := json.NewDecoder(file).Decode(&settings); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	return settings.ProvincialCenters, nil
}
