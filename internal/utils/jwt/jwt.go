package jwt

import (
	domainerrors "capecom-pm/internal/domain/error"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenType string

const (
	TokenTypeUser         TokenType = "user"
	TokenTypeAdmin        TokenType = "admin"
	TokenTypeRefresh      TokenType = "refresh"
	TokenTypeAdminRefresh TokenType = "admin_refresh"
)

type Claims struct {
	UserID string `json:"uid"`
	jwt.RegisteredClaims
}

type Config struct {
	UserSecret         string
	AdminSecret        string
	ExpireHours        int
	RefreshExpireHours int
}

type JWTManager struct {
	config Config
}

func NewJWTManager(userSecret, adminSecret string, expireHours, refreshExpireHours int) *JWTManager {
	return &JWTManager{
		config: Config{
			UserSecret:         userSecret,
			AdminSecret:        adminSecret,
			ExpireHours:        expireHours,
			RefreshExpireHours: refreshExpireHours,
		},
	}
}

// CreateToken creates a JWT token based on token type
func (j *JWTManager) CreateToken(userID string, tokenType TokenType) (string, error) {
	var expirationTime time.Time

	// Use different expiration for refresh tokens
	if tokenType == TokenTypeRefresh || tokenType == TokenTypeAdminRefresh {
		expirationTime = time.Now().Add(time.Duration(j.config.RefreshExpireHours) * time.Hour)
	} else {
		expirationTime = time.Now().Add(time.Duration(j.config.ExpireHours) * time.Hour)
	}

	jti := uuid.New().String()

	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        jti,
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Use different secret based on token type
	secret := j.getSecret(tokenType)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates and parses a JWT token
func (j *JWTManager) ValidateToken(tokenString string, tokenType TokenType) (*Claims, error) {
	claims := &Claims{}

	secret := j.getSecret(tokenType)

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, domainerrors.ErrInvalidSigningMethod
		}
		return []byte(secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, domainerrors.ErrTokenExpired
		}
		return nil, domainerrors.ErrInvalidToken
	}

	if !token.Valid {
		return nil, domainerrors.ErrInvalidToken
	}

	return claims, nil
}

// ValidateAnyToken validates token without checking type (useful for middleware)
func (j *JWTManager) ValidateAnyToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	// Try user secret first, then admin secret
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, domainerrors.ErrInvalidSigningMethod
		}
		return []byte(j.config.UserSecret), nil
	})

	// If failed with user secret, try admin secret
	if err != nil || !token.Valid {
		token, err = jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, domainerrors.ErrInvalidSigningMethod
			}
			return []byte(j.config.AdminSecret), nil
		})
	}

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, domainerrors.ErrTokenExpired
		}
		return nil, domainerrors.ErrInvalidToken
	}

	if !token.Valid {
		return nil, domainerrors.ErrInvalidToken
	}

	return claims, nil
}

func (j *JWTManager) getSecret(tokenType TokenType) string {
	if tokenType == TokenTypeAdmin || tokenType == TokenTypeAdminRefresh {
		return j.config.AdminSecret
	}
	return j.config.UserSecret
}
