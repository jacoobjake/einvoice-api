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
	"github.com/gofrs/uuid/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jacoobjake/einvoice-api/config"
	"github.com/jacoobjake/einvoice-api/internal/database/enums"
	"github.com/jacoobjake/einvoice-api/internal/database/models"
	"github.com/jacoobjake/einvoice-api/internal/repositories"
	"github.com/jacoobjake/einvoice-api/pkg"
	"github.com/jacoobjake/einvoice-api/pkg/redisclient"
)

type AuthService struct {
	authRepo      *repositories.AuthTokenRepository
	userRepo      *repositories.UserRepository
	config        *config.Config
	signingMethod jwt.SigningMethod
	rdb           *redisclient.RedisClient
	revokedPrefix string
}

type AuthClaims struct {
	UserID    int64
	Email     string
	SessionID uuid.UUID
	jwt.RegisteredClaims
}

func (s *AuthService) getRevokedTokenKey(token string) string {
	return fmt.Sprintf("%s%s", s.revokedPrefix, token)
}

func (s *AuthService) hashRefreshToken(token string) (string, error) {
	encrypted := hmac.New(sha256.New, []byte(s.config.AuthConfig.RefreshTokenSecret))
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

func (s *AuthService) invalidateActiveRefreshTokens(ctx context.Context, sessionId uuid.UUID) error {
	err := s.authRepo.InvalidateActiveTokensBySessionID(ctx, sessionId, enums.AuthTokenTypesRefresh)
	return err
}

func (s *AuthService) generateRefreshToken(ctx context.Context, user models.User, sessionId uuid.UUID) (string, error) {
	err := s.invalidateActiveRefreshTokens(ctx, sessionId)
	authConfig := s.config.AuthConfig

	if err != nil {
		return "", err
	}

	refreshToken, err := pkg.GenerateRandomString(32)
	duration := time.Duration(authConfig.RefreshExpirationMin) * time.Minute

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
		UserID:    omit.From(user.ID),
		ExpireAt:  omitnull.From(time.Now().Add(duration)),
		Type:      omit.From(enums.AuthTokenTypesRefresh),
		Token:     omit.From(hashed),
		SessionID: omitnull.From(sessionId),
	})

	if err != nil {
		return "", err
	}

	return refreshToken, nil
}

func (s *AuthService) generateToken(ctx context.Context, user models.User) (token string, refreshToken string, err error) {
	var t *jwt.Token
	authConfig := s.config.AuthConfig
	key := []byte(authConfig.JWTSecret)
	sessionId := uuid.Must(uuid.NewV4())

	claims := AuthClaims{
		user.ID,
		user.Email,
		sessionId,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(authConfig.TokenExpirationMin) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    s.config.AppName,
		},
	}

	t = jwt.NewWithClaims(s.signingMethod, claims)

	signed, err := t.SignedString(key)

	if err != nil {
		return "", "", err
	}

	refreshToken, err = s.generateRefreshToken(ctx, user, sessionId)

	if err != nil {
		return "", "", err
	}

	return signed, refreshToken, nil
}

func (s *AuthService) parseToken(_ context.Context, token string) (claims jwt.Claims, err error) {
	parsed, err := jwt.ParseWithClaims(token, &AuthClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(s.config.AuthConfig.JWTSecret), nil
	}, jwt.WithValidMethods([]string{s.signingMethod.Alg()}))

	if err != nil {
		return nil, err
	}

	claims, ok := parsed.Claims.(*AuthClaims)

	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}

	return claims, nil
}

func (s *AuthService) verifyJWTToken(ctx context.Context, token string) (*AuthClaims, error) {
	claims, err := s.parseToken(ctx, token)

	if err != nil {
		return nil, err
	}

	authClaims, ok := claims.(*AuthClaims)

	if !ok {
		return nil, fmt.Errorf("invalid claim type")
	}

	// Check expiry
	exp, err := authClaims.GetExpirationTime()

	if err != nil {
		return nil, err
	}

	if time.Now().After(exp.Time) {
		return nil, fmt.Errorf("token expired")
	}

	// Check not before
	nbf, err := authClaims.GetNotBefore()

	if err != nil {
		return nil, err
	}

	if time.Now().Before(nbf.Time) {
		return nil, fmt.Errorf("token is not valid yet")
	}

	// Check if token is revoked
	key := s.getRevokedTokenKey(token)
	revoked, err := s.rdb.Exists(ctx, key)

	if err != nil {
		return nil, fmt.Errorf("error reading key: %w", err)
	}

	if revoked {
		return nil, fmt.Errorf("token revoked")
	}

	return authClaims, nil
}

func (s *AuthService) isActiveUser(user *models.User) bool {
	return user.DeletedAt.IsNull() && user.Status == enums.UserStatusesActive
}

func (s *AuthService) Token(ctx context.Context, email string, pw string) (rawToken string, refreshToken string, err error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", "", fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return "", "", fmt.Errorf("user not found")
	}
	if err := pkg.ComparePassword([]byte(user.Password), []byte(pw)); err != nil {
		return "", "", fmt.Errorf("password mismatch")
	}
	// TODO:  Check failed login attempts
	rawToken, refreshToken, err = s.generateToken(ctx, *user)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate token: %w", err)
	}
	return rawToken, refreshToken, nil
}

func (s *AuthService) RevokeToken(ctx context.Context, token string) error {
	key := s.getRevokedTokenKey(token)

	exists, err := s.rdb.Exists(ctx, key)

	if err != nil {
		return fmt.Errorf("failed to read key: %w", err)
	}

	// Only revoke if not already revoked
	if exists {
		return nil
	}

	err = s.rdb.Set(ctx, key, true, time.Duration(15*time.Minute))

	if err != nil {
		return fmt.Errorf("failed to write key: %w", err)
	}

	claims, err := s.parseToken(ctx, token)

	if err != nil {
		return fmt.Errorf("error parsing token: %w", err)
	}

	authClaims := claims.(*AuthClaims)

	// Invalidate refresh token
	err = s.invalidateActiveRefreshTokens(ctx, authClaims.SessionID)

	if err != nil {
		return fmt.Errorf("error invalidating refresh token: %w", err)
	}

	return nil
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

func (s *AuthService) VerifyToken(ctx context.Context, token string) (*models.User, error) {
	claims, err := s.verifyJWTToken(ctx, token)

	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindById(ctx, claims.UserID)

	if err != nil {
		return nil, fmt.Errorf("invalid claim user id: %d", claims.UserID)
	}

	if !s.isActiveUser(user) {
		return nil, fmt.Errorf("user account inactive")
	}

	return user, nil
}

func NewAuthService(
	authRepo *repositories.AuthTokenRepository,
	userRepo *repositories.UserRepository,
	config *config.Config,
	rdb *redisclient.RedisClient,
) *AuthService {
	return &AuthService{
		authRepo:      authRepo,
		userRepo:      userRepo,
		config:        config,
		signingMethod: jwt.SigningMethodHS256,
		rdb:           rdb,
		revokedPrefix: "revoked:",
	}
}
