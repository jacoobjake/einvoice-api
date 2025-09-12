package services

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jacoobjake/einvoice-api/config/auth"
	"github.com/jacoobjake/einvoice-api/internal/database/enums"
	"github.com/jacoobjake/einvoice-api/internal/database/models"
	"github.com/jacoobjake/einvoice-api/internal/repositories"
	"github.com/jacoobjake/einvoice-api/pkg"
)

type AuthService struct {
	authRepo   *repositories.AuthTokenRepository
	userRepo   *repositories.UserRepository
	authConfig *auth.AuthConfig
}

func (s *AuthService) hashRefreshToken(token string) (string, error) {
	encrypted := hmac.New(sha256.New, []byte(s.authConfig.RefreshTokenSecret))
	_, err := encrypted.Write([]byte(token))

	if err != nil {
		return "", err
	}

	return hex.EncodeToString(encrypted.Sum(nil)), nil
}

func (s *AuthService) validateRefreshToken(refreshToken *models.AuthToken, plainToken string) bool {
	if refreshToken.Type != enums.AuthTokenTypesRefresh {
		return false
	}

	expireAt, isset := refreshToken.ExpireAt.Get()

	if !isset || !expireAt.After(time.Now()) {
		return false
	}

	hashedToken := refreshToken.Token
	encrypted, err := s.hashRefreshToken(plainToken)

	if err != nil {
		return false
	}

	return hashedToken == encrypted
}

func (s *AuthService) invalidateActiveRefreshTokens(ctx context.Context, userID int64) error {
	err := s.authRepo.InvalidateActiveTokensByUserID(ctx, userID, enums.AuthTokenTypesRefresh)
	return err
}

func (s *AuthService) generateRefreshToken(ctx context.Context, user models.User) (string, error) {
	err := s.invalidateActiveRefreshTokens(ctx, user.ID)

	if err != nil {
		return "", fmt.Errorf("could not invalidate existing tokens: %w", err)
	}

	refreshToken, err := pkg.GenerateRandomString(32)
	duration := time.Duration(s.authConfig.RefreshExpirationMin) * time.Minute

	if err != nil {
		return "", err
	}

	// Store encrypted version in DB
	hashed, err := s.hashRefreshToken(refreshToken)

	if err != nil {
		return "", err
	}

	// Store refresh token in DB
	_, err = s.authRepo.Create(ctx, &models.AuthTokenSetter{
		UserID:   omit.From(user.ID),
		ExpireAt: omitnull.From(time.Now().Add(duration)),
		Type:     omit.From(enums.AuthTokenTypesRefresh),
		Token:    omit.From(hashed),
	})

	if err != nil {
		return "", err
	}

	return refreshToken, nil
}

func (s *AuthService) generateToken(ctx context.Context, user models.User) (token string, refreshToken string, err error) {
	var t *jwt.Token
	key := []byte(s.authConfig.JWTSecret)

	t = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"iat":     time.Now().Unix(),
		"exp":     time.Now().Add(time.Duration(s.authConfig.TokenExpirationMin) * time.Minute).Unix(),
		"nbf":     time.Now().Unix(),
	})

	signed, err := t.SignedString(key)

	if err != nil {
		return "", "", err
	}

	refreshToken, err = s.generateRefreshToken(ctx, user)

	if err != nil {
		return "", "", err
	}

	return signed, refreshToken, nil
}

func (s *AuthService) Login(ctx context.Context, email string, pw string) (rawToken string, refreshToken string, err error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", "", fmt.Errorf("failed to find user")
	}
	if user == nil {
		return "", "", fmt.Errorf("user not found")
	}
	if err := pkg.ComparePassword([]byte(user.Password), []byte(pw)); err != nil {
		return "", "", fmt.Errorf("password mismatch")
	}
	rawToken, refreshToken, err = s.generateToken(ctx, *user)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate token")
	}
	return rawToken, refreshToken, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, email string, refreshToken string) (rawToken string, newRefreshToken string, err error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", "", fmt.Errorf("failed to find user")
	}
	if user == nil {
		return "", "", fmt.Errorf("user not found")
	}
	storedToken, err := s.authRepo.FindTokenByUserIdAndType(ctx, user.ID, enums.AuthTokenTypesRefresh)
	if err != nil {
		return "", "", fmt.Errorf("failed to find refresh token")
	}
	if storedToken == nil || storedToken.UserID != user.ID {
		return "", "", fmt.Errorf("refresh token not found or does not belong to user")
	}
	if !s.validateRefreshToken(storedToken, refreshToken) {
		return "", "", fmt.Errorf("refresh token validation failed")
	}
	rawToken, newRefreshToken, err = s.generateToken(ctx, *user)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate new token")
	}
	return rawToken, newRefreshToken, nil
}

func NewAuthService(
	authRepo *repositories.AuthTokenRepository,
	userRepo *repositories.UserRepository,
	authConfig *auth.AuthConfig,
) *AuthService {
	return &AuthService{authRepo: authRepo, userRepo: userRepo, authConfig: authConfig}
}
