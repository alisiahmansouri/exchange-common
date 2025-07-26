package jwt

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken   = errors.New("invalid token")
	ErrTokenRevoked   = errors.New("token revoked")
	ErrRefreshExpired = errors.New("refresh token expired")
)

// Service handles JWT creation, validation, revocation, and refresh token rotation.
type Service struct {
	secretKey   string
	tokenTTL    time.Duration
	refreshTTL  time.Duration
	redisClient *redis.Client
	redisCtx    context.Context
}

// New creates a new Service instance.
func New(secret string, tokenTTL, refreshTTL time.Duration, rdb *redis.Client) *Service {
	return &Service{
		secretKey:   secret,
		tokenTTL:    tokenTTL,
		refreshTTL:  refreshTTL,
		redisClient: rdb,
		redisCtx:    context.Background(),
	}
}

// Claims represents the JWT claims, including user ID and standard registered claims.
type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func (s *Service) GenerateToken(userID uuid.UUID) (accessToken string, refreshToken string, err error) {
	accessToken, err = s.GenerateAccessToken(userID)
	if err != nil {
		return "", "", err
	}
	refreshToken, err = s.GenerateRefreshToken(userID.String())
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

// GenerateAccessToken creates a new JWT access token for the given user ID.
// The token includes an expiration time and a unique ID for revocation tracking.
func (s *Service) GenerateAccessToken(userID uuid.UUID) (string, error) {
	claims := Claims{
		UserID: userID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.tokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        uuid.NewString(), // Unique ID for token revocation
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secretKey))
}

// GenerateRefreshToken creates a new JWT refresh token with longer expiration.
func (s *Service) GenerateRefreshToken(userID string) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.refreshTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        uuid.NewString(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secretKey))
}

// ValidateToken parses and validates the token string, returning claims if valid.
// It also checks if the token is revoked in Redis blacklist.
func (s *Service) ValidateToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(s.secretKey), nil
	})

	if err != nil {
		// Check if the error is due to token expiration
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrRefreshExpired
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	// Check if token is revoked in Redis blacklist
	revoked, err := s.isTokenRevoked(claims.ID)
	if err != nil {
		return nil, err
	}
	if revoked {
		return nil, ErrTokenRevoked
	}

	return claims, nil
}

// RevokeToken adds the token's unique ID (jti) to Redis blacklist with TTL until token expiration.
func (s *Service) RevokeToken(tokenStr string) error {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(s.secretKey), nil
	})
	if err != nil {
		return ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return ErrInvalidToken
	}

	exp := claims.ExpiresAt.Time
	ttl := time.Until(exp)
	if ttl <= 0 {
		// Token already expired, no need to blacklist
		return nil
	}

	// Store token ID in Redis with expiration to blacklist it
	return s.redisClient.Set(s.redisCtx, claims.ID, "revoked", ttl).Err()
}

// isTokenRevoked checks Redis to see if the token ID is blacklisted.
func (s *Service) isTokenRevoked(jti string) (bool, error) {
	res, err := s.redisClient.Exists(s.redisCtx, jti).Result()
	if err != nil {
		return false, err
	}
	return res == 1, nil
}

// RotateRefreshToken revokes the old refresh token and issues a new one for the same user.
func (s *Service) RotateRefreshToken(oldToken string) (string, error) {
	claims, err := s.ValidateToken(oldToken)
	if err != nil {
		return "", err
	}

	// Determine if token is a refresh token by checking its expiration duration
	if time.Until(claims.ExpiresAt.Time) > s.tokenTTL {
		// Revoke old refresh token
		err = s.RevokeToken(oldToken)
		if err != nil {
			return "", err
		}

		// Generate new refresh token
		return s.GenerateRefreshToken(claims.UserID)
	}
	return "", errors.New("token is not a refresh token or expired")
}

func (s *Service) ValidateRefreshToken(tokenStr string) (*Claims, error) {
	claims, err := s.ValidateToken(tokenStr)
	if err != nil {
		return nil, err
	}

	// چک کن که این توکن واقعاً Refresh Token هست یا نه
	// Refresh Token باید زمان انقضاش بیشتر از access token باشه
	if time.Until(claims.ExpiresAt.Time) <= s.tokenTTL {
		return nil, errors.New("token is not a refresh token or expired")
	}

	return claims, nil
}
