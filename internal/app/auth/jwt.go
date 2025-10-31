package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// JWTClaims - структура для JWT токена
type JWTClaims struct {
	UserUUID string `json:"user_uuid"`
	Role     string `json:"role"`
	Type     string `json:"type"` // "access" или "refresh"
	jwt.RegisteredClaims
}

// JWTService - сервис для работы с JWT
type JWTService struct {
	secretKey              string
	accessTokenExpiration  time.Duration
	refreshTokenExpiration time.Duration
}

// NewJWTService - создание нового JWT сервиса
func NewJWTService(secretKey string, accessExpMin, refreshExpDays int) *JWTService {
	return &JWTService{
		secretKey:              secretKey,
		accessTokenExpiration:  time.Duration(accessExpMin) * time.Minute,
		refreshTokenExpiration: time.Duration(refreshExpDays) * 24 * time.Hour,
	}
}

// GenerateAccessToken - генерация access токена
func (j *JWTService) GenerateAccessToken(userUUID, role string) (string, error) {
	claims := JWTClaims{
		UserUUID: userUUID,
		Role:     role,
		Type:     "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.accessTokenExpiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

// GenerateRefreshToken - генерация refresh токена
func (j *JWTService) GenerateRefreshToken(userUUID, role string) (string, error) {
	claims := JWTClaims{
		UserUUID: userUUID,
		Role:     role,
		Type:     "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.refreshTokenExpiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

// ValidateToken - валидация токена
func (j *JWTService) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(j.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// RefreshTokenPair - обновление пары токенов
func (j *JWTService) RefreshTokenPair(refreshToken string) (string, string, error) {
	claims, err := j.ValidateToken(refreshToken)
	if err != nil {
		return "", "", err
	}

	if claims.Type != "refresh" {
		return "", "", errors.New("invalid token type")
	}

	accessToken, err := j.GenerateAccessToken(claims.UserUUID, claims.Role)
	if err != nil {
		return "", "", err
	}

	newRefreshToken, err := j.GenerateRefreshToken(claims.UserUUID, claims.Role)
	if err != nil {
		return "", "", err
	}

	return accessToken, newRefreshToken, nil
}

// GetTokenExpiration - получение времени истечения токена
func (j *JWTService) GetTokenExpiration(tokenString string) (time.Duration, error) {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return 0, err
	}

	if claims.ExpiresAt == nil {
		return 0, errors.New("token has no expiration")
	}

	return time.Until(claims.ExpiresAt.Time), nil
}

