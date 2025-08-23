package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var (
	BaseURL  string
	USERNAME string
	PASSWORD string
)

func LoadConfig() {
	// Load .env jika ada
	if err := godotenv.Load(); err != nil {
		fmt.Println("[WARN] .env tidak ditemukan, gunakan env sistem")
	}

	BaseURL = os.Getenv("BASE_URL")
	USERNAME = os.Getenv("USER_SIAKAD")
	PASSWORD = os.Getenv("PASSWORD_SIAKAD")
	if BaseURL == "" {
		fmt.Println("[ERROR] BASE_URL tidak ditemukan di .env atau env sistem")
		os.Exit(1)
	}
	if USERNAME == "" {
		fmt.Println("[ERROR] USERNAME tidak ditemukan di .env atau env sistem")
		os.Exit(1)
	}
	if PASSWORD == "" {
		fmt.Println("[ERROR] PASSWORD tidak ditemukan di .env atau env sistem")
		os.Exit(1)
	}

}
