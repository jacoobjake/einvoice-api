package env

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() string {
	// Load layered .env files
	env := os.Getenv("FOO_ENV")
	if env == "" {
		env = "development"
	}

	_ = godotenv.Load(".env." + env + ".local")

	if env != "test" {
		_ = godotenv.Load(".env.local")
	}
	_ = godotenv.Load(".env." + env)
	_ = godotenv.Load() // fallback .env

	return env
}

func GetEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func GetEnvAsInt(key string, fallback int) int {
	if value, exists := os.LookupEnv(key); exists {
		var intVal int
		_, err := fmt.Sscanf(value, "%d", &intVal)
		if err == nil {
			return intVal
		}
	}
	return fallback
}
