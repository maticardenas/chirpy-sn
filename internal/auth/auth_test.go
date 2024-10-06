package auth

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
)

func TestMakeJWT(t *testing.T) {
	// slice with test cases
	testCases := []struct {
		userId      uuid.UUID
		tokenSecret string
		expiresIn   int
	}{
		{
			userId:      uuid.New(),
			tokenSecret: "secret",
			expiresIn:   3600,
		},
	}

	for _, tc := range testCases {
		token, err := MakeJWT(tc.userId, tc.tokenSecret, tc.expiresIn)
		if err != nil {
			t.Errorf("Expected no error, but got %v", err)
		}
		fmt.Printf("Token: %v\n", token)
		fmt.Printf("Token length: %v\n", len(token))
	}
}

func TestValidateJWT(t *testing.T) {
	tokenSecret := "secret"
	token, err := MakeJWT(uuid.New(), tokenSecret, 3600)

	if err != nil {
		t.Errorf("Expected no erro creating JWT, but got %v", err)
	}

	userId, err := ValidateJWT(token, tokenSecret)

	if err != nil {
		t.Errorf("Expected no error validating JWT, but got %v", err)
	}

	fmt.Printf("User ID: %v\n", userId)
}

func TestGetBearerToken(t *testing.T) {
	// slice with test cases
	testCases := []struct {
		headers  http.Header
		expected string
	}{
		{
			headers: http.Header{
				"Authorization": []string{"Bearer token"},
			},
			expected: "token",
		},
		{
			headers: http.Header{
				"Authorization": []string{"Auguri token"},
			},
			expected: "",
		},
		{
			headers: http.Header{
				"Authorization": []string{""},
			},
			expected: "",
		},
		{
			headers: http.Header{
				"Authorization": []string{"token"},
			},
			expected: "",
		},
	}

	for _, tc := range testCases {
		bearerToken, err := GetBearerToken(tc.headers)

		if err != nil {
			if tc.expected != "" {
				t.Errorf("Expected no error, but got %v", err)
			}
		}

		if bearerToken != tc.expected {
			t.Errorf("Expected %v, but got %v", tc.expected, bearerToken)
		}
	}
}
