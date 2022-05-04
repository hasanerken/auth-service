package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

func Get(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading environment file.")
	}
	return os.Getenv(key)
}
