package main

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var secretKey = []byte("your-secret-key") // Your secret key to sign the JWT

// GenerateJWT generates a JWT for a user with a userID
func GenerateJWT(userID int) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,                                // The user ID is stored in the "sub" claim
		"exp": time.Now().Add(time.Hour * 24).Unix(), // Set the token expiration time
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// ParseToken parses and validates the JWT token, returning the user ID
func ParseToken(tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secretKey, nil
	})

	if err != nil {
		return 0, err
	}

	// Extract the user ID from the claims if the token is valid
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := int(claims["sub"].(float64)) // Extract the "sub" claim (user ID)
		return userID, nil
	}

	return 0, errors.New("invalid token")
}
