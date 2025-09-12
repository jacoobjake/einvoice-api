package auth

import (
	"github.com/jacoobjake/einvoice-api/pkg/env"
)

type AuthConfig struct {
	JWTSecret            string
	RefreshTokenSecret   string
	TokenExpirationMin   int
	RefreshExpirationMin int
}

func LoadAuthConfig() *AuthConfig {
	return &AuthConfig{
		JWTSecret:            env.GetEnv("JWT_SECRET", "default_jwt_secret"),
		RefreshTokenSecret:   env.GetEnv("REFRESH_TOKEN_SECRET", "default_refresh_token_secret"),
		TokenExpirationMin:   env.GetEnvAsInt("TOKEN_EXPIRATION_MIN", 15),
		RefreshExpirationMin: env.GetEnvAsInt("REFRESH_EXPIRATION_MIN", 24*60),
	}
}
