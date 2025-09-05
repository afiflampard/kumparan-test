package config

import (
	"log"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

func LoadConfig() Config {
	_, filename, _, _ := runtime.Caller(0)

	rootPath := filepath.Join(filepath.Dir(filename), "..")

	envPath := filepath.Join(rootPath, ".env")
	if err := godotenv.Load(envPath); err != nil {
		log.Printf("⚠️ Could not load .env file: %v", err)
	}

	var AppConfig Config
	if err := envconfig.Process("", &AppConfig); err != nil {
		log.Fatalf("Failed to load env config: %v", err)
	}
	return AppConfig
}
