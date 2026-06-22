package security

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type TokenClaims struct {
	UserID uuid.UUID `json:"user_id"`

	jwt.RegisteredClaims
}

func GenerateToken(
	userID uuid.UUID,
	privateKeyPath string,
	tokenTTL time.Duration,
) (string, error) {
	privateKeyBytes, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return "", fmt.Errorf("read private key: %w", err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
		return "", fmt.Errorf("parse private key: %w", err)
	}

	now := time.Now()

	claims := TokenClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(tokenTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodRS256,
		claims,
	)

	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}

	return tokenString, nil
}

func ParseToken(
	tokenString string,
	publicKeyPath string,
) (*TokenClaims, error) {
	publicKeyBytes, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("read public key: %w", err)
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("parse public key: %w", err)
	}

	claims := &TokenClaims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			if token.Method.Alg() != jwt.SigningMethodRS256.Alg() {
				return nil, fmt.Errorf(
					"unexpected signing method: %s",
					token.Method.Alg(),
				)
			}

			return publicKey, nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("token is invalid")
	}

	return claims, nil
}
