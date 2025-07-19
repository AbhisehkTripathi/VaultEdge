package config

import (
	"log"
	"os"
)

// LoadEnv loads environment variables from .env file if present
func LoadEnv() {
	if _, err := os.Stat(".env"); err == nil {
		if err := loadDotEnv(); err != nil {
			log.Printf("Could not load .env: %v", err)
		}
	}
}

// loadDotEnv loads .env (helper for LoadEnv)
func loadDotEnv() error {
	// Use os package to load .env manually (if you want to avoid external deps)
	// Or use github.com/joho/godotenv for robust loading
	return nil
}
