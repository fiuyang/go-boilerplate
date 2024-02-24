package utils

import (
	"fmt"
	"sync"
	"time"

	"github.com/golang-jwt/jwt"
)

var (
	blacklistMutex sync.Mutex
	blacklist      = make(map[string]bool)
)

func GenerateToken(ttl time.Duration, payload interface{}, secretJWTKey string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	now := time.Now().UTC()
	claims := token.Claims.(jwt.MapClaims)

	claims["sub"] = payload
	claims["exp"] = now.Add(ttl).Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	tokenString, err := token.SignedString([]byte(secretJWTKey))

	if err != nil {
		return "", fmt.Errorf("generating JWT Token failed: %w", err)
	}

	return tokenString, nil
}

func RefreshToken(oldToken string, ttl time.Duration, secretJWTKey string) (string, error) {
	// Parse token
	oldParsedToken, err := jwt.Parse(oldToken, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected method: %s", jwtToken.Header["alg"])
		}

		return []byte(secretJWTKey), nil
	})
	if err != nil {
		return "", fmt.Errorf("invalidate token: %w", err)
	}

	// Check if the token is valid
	if !oldParsedToken.Valid {
		return "", fmt.Errorf("old token is not valid")
	}

	// Extract claims from the old token
	claims, ok := oldParsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("failed to extract claims from old token")
	}

	// Generate a new token with the same payload but new expiration time
	now := time.Now().UTC()
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": claims["sub"],
		"exp": now.Add(ttl).Unix(),
		"iat": now.Unix(),
		"nbf": now.Unix(),
	})

	// Sign the token with the secret key
	newTokenString, err := newToken.SignedString([]byte(secretJWTKey))
	if err != nil {
		return "", fmt.Errorf("failed to generate new token: %w", err)
	}

	return newTokenString, nil
}

func ValidateToken(token string, signedJWTKey string) (interface{}, error) {
	tok, err := jwt.Parse(token, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected method: %s", jwtToken.Header["alg"])
		}

		return []byte(signedJWTKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalidate token: %w", err)
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok || !tok.Valid {
		return nil, fmt.Errorf("invalid token claim")
	}

	return claims["sub"], nil
}

func AddToBlacklist(token string) error {
	blacklistMutex.Lock()
	defer blacklistMutex.Unlock()

	blacklist[token] = true
	return nil
}

func IsTokenBlacklisted(token string) bool {
	blacklistMutex.Lock()
	defer blacklistMutex.Unlock()
	return blacklist[token]
}
