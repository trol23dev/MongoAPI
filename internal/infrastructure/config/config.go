package config

import (
	"os"
	"path/filepath"

	"mongodb-api/internal/infrastructure/logger"

	"github.com/joho/godotenv"
)

type Config struct {
	Port     string
	MongoURI string
}

func Load() *Config {
	paths := []string{
		".env",
		"../.env",
		"../../.env",
	}

	if execPath, err := os.Executable(); err == nil {
		execDir := filepath.Dir(execPath)
		paths = append(paths, filepath.Join(execDir, ".env"), filepath.Join(execDir, "..", ".env"))
	}

	loaded := false
	for _, p := range paths {
		if err := godotenv.Load(p); err == nil {
			loaded = true
			break
		}
	}

	mongoURI := os.Getenv("MONGO_URI")
	if !loaded && mongoURI == "" {
		logger.New(logger.Warning).Print("No se pudo cargar el archivo .env desde las rutas estándar ni se encontró MONGO_URI en el sistema.")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	return &Config{
		Port:     port,
		MongoURI: mongoURI,
	}
}
