package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config holds the application configuration
type Config struct {
	BaseURL  string
	Username string
	Password string
}

// LoadConfig loads configuration from environment variables or .env file
func LoadConfig() (*Config, error) {
	// Load .env jika ada
	if err := godotenv.Load(); err != nil {
		fmt.Println("[WARN] .env tidak ditemukan, gunakan env sistem")
	}

	config := &Config{
		BaseURL:  os.Getenv("BASE_URL"),
		Username: os.Getenv("USER_SIAKAD"),
		Password: os.Getenv("PASSWORD_SIAKAD"),
	}

	if config.BaseURL == "" {
		return nil, fmt.Errorf("BASE_URL tidak ditemukan di .env atau env sistem")
	}
	if config.Username == "" {
		return nil, fmt.Errorf("USER_SIAKAD tidak ditemukan di .env atau env sistem")
	}
	if config.Password == "" {
		return nil, fmt.Errorf("PASSWORD_SIAKAD tidak ditemukan di .env atau env sistem")
	}

	return config, nil
}
