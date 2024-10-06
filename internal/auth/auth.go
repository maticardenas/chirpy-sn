package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userId uuid.UUID, tokenSecret string, expiresIn int) (string, error) {
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiresIn) * time.Second).UTC()),
		Subject:   userId.String(),
	}).SignedString([]byte(tokenSecret))

	if err != nil {
		fmt.Println("Error creating JWT token:", err)
		return "", err
	}

	return token, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})

	if err != nil {
		fmt.Println("Error parsing JWT token:", err)
		return uuid.UUID{}, err
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)

	fmt.Printf("Claims: %v\n", claims)

	if !ok || !token.Valid {
		return uuid.UUID{}, fmt.Errorf("invalid token")
	}

	userId, err := claims.GetSubject()
	if err != nil {
		fmt.Println("Error getting subject from JWT token:", err)
		return uuid.UUID{}, err
	}
	return uuid.Parse(userId)
}

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")

	fmt.Printf("Authorization header: %v\n", authHeader)
	if authHeader == "" {
		return "", fmt.Errorf("No Authorization header provided")
	}

	splitAuthHeader := strings.Split(authHeader, " ")

	fmt.Printf("Split auth header: %v\n", splitAuthHeader)

	if len(splitAuthHeader) != 2 || strings.ToLower(splitAuthHeader[0]) != "bearer" {
		return "", fmt.Errorf("Invalid Authorization header")
	}

	return splitAuthHeader[1], nil
}
