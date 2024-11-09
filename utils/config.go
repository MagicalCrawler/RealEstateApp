package utils

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func GetConfig(key string) string {
	return os.Getenv(key)
}

func LoadEnvFile() {
	err := godotenv.Load(".env.example")
	if err != nil {
		fmt.Println(err)
		panic("Error loading .env.example file")
	}
}
