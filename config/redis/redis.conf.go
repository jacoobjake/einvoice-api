package redis

import "github.com/jacoobjake/einvoice-api/pkg/env"

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

func LoadRedisConfig() *RedisConfig {
	return &RedisConfig{
		Addr:     env.GetEnv("REDIS_ADDR", "localhost:6379"),
		Password: env.GetEnv("REDIS_PASSWORD", ""),
		DB:       env.GetEnvAsInt("REDIS_DB", 0),
	}
}
