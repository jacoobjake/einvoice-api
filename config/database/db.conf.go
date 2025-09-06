package database

import (
	"fmt"

	"github.com/jacoobjake/einvoice-api/pkg/env"
)

type DBConfig struct {
	Driver   string
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
	TimeZone string
}

func LoadDBConfig() *DBConfig {
	return &DBConfig{
		Driver:   env.GetEnv("DB_DRIVER", "postgres"),
		Host:     env.GetEnv("DB_HOST", "localhost"),
		Port:     env.GetEnv("DB_PORT", "5432"),
		User:     env.GetEnv("DB_USER", "root"),
		Password: env.GetEnv("DB_PASSWORD", "password"),
		DBName:   env.GetEnv("DB_NAME", "einvoice"),
		SSLMode:  env.GetEnv("DB_SSLMODE", "disable"),
		TimeZone: env.GetEnv("DB_TIMEZONE", "TimeZone=Asia/Kuala_Lumpur"),
	}
}

func (cfg *DBConfig) DSN() string {
	// Format DSB base
	switch cfg.Driver {
	case "postgres":
		return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s %s",
			cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port, cfg.SSLMode, cfg.TimeZone)
	default:
		return ""
	}
}

func (cfg *DBConfig) ConnectionString() string {
	// Driver string format
	switch cfg.Driver {
	case "postgres":
		return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s&%s",
			cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode, cfg.TimeZone)
	default:
		return ""
	}
}
