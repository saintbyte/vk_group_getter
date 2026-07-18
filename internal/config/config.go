package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	VKToken string
	DBPath  string
	APIKey  string
	Host    string
	Port    string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &Config{
		VKToken: os.Getenv("VK_TOKEN"),
		DBPath:  os.Getenv("DB_PATH"),
		APIKey:  os.Getenv("API_KEY"),
		Host:    os.Getenv("HOST"),
		Port:    port,
	}
}
