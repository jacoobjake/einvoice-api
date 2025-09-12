package config

import (
	"github.com/jacoobjake/einvoice-api/config/auth"
	"github.com/jacoobjake/einvoice-api/config/database"
	pkgEnv "github.com/jacoobjake/einvoice-api/pkg/env"
)

type Config struct {
	AppName    string
	Port       string
	Env        string
	DBConfig   *database.DBConfig
	AuthConfig *auth.AuthConfig
}

func Load() *Config {
	// Load layered .env files
	env := pkgEnv.LoadEnv()
	DBConfig := database.LoadDBConfig()
	AuthConfig := auth.LoadAuthConfig()

	cfg := &Config{
		AppName:    pkgEnv.GetEnv("APP_NAME", "MyApp"),
		Port:       pkgEnv.GetEnv("PORT", "8080"),
		DBConfig:   DBConfig,
		AuthConfig: AuthConfig,
		Env:        env,
	}

	return cfg
}
