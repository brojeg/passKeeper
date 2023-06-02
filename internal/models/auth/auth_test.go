package models

import (
	"context"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func TestValidateToken(t *testing.T) {
	InitJWTPassword("testpassword", 15)
	expectedUserID := uint(1)

	tk := &Token{UserID: expectedUserID}
	expirationTime := time.Now().Add(time.Duration(settings.expirationTime) * time.Minute)
	tk.ExpiresAt = expirationTime.Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tk)
	tokenString, err := token.SignedString([]byte(settings.jwtPassword))
	if err != nil {
		t.Fatalf("unable to sign token: %v", err)
	}

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	req.Header.Set("Authorization", tokenString)

	rr := ValidateToken(req)
	if rr.ServerCode != 200 {
		t.Errorf("handler returned unexpected status code: got %v want %v", rr.ServerCode, 200)
	}

	if strconv.FormatUint(uint64(rr.Message.(uint)), 10) != strconv.Itoa(int(expectedUserID)) {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Message, strconv.Itoa(int(expectedUserID)))
	}
}

func TestIsPasswordsEqual(t *testing.T) {
	// Given
	password := "testpassword"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	testCases := []struct {
		name     string
		existing string
		new      string
		expected bool
	}{
		{
			name:     "Equal passwords",
			existing: string(hashedPassword),
			new:      password,
			expected: true,
		},
		{
			name:     "Not equal passwords",
			existing: string(hashedPassword),
			new:      "wrongpassword",
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// When
			result := IsPasswordsEqual(tc.existing, tc.new)

			// Then
			if result != tc.expected {
				t.Errorf("Expected %v, got %v", tc.expected, result)
			}
		})
	}
}
func TestEncryptPassword(t *testing.T) {
	// Given
	password := "testpassword"

	// When
	hashedPassword := EncryptPassword(password)

	// Then
	if hashedPassword == password {
		t.Errorf("Password was not hashed")
	}

	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		t.Errorf("Password hashing is not correct")
	}
}
func TestGenerateToken(t *testing.T) {
	// Given
	InitJWTPassword("testpassword", 5)
	userID := uint(1)

	// When
	tokenString := GenerateToken(userID)

	// Then
	token, err := jwt.ParseWithClaims(tokenString, &Token{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(settings.jwtPassword), nil
	})

	if err != nil {
		t.Errorf("Token could not be parsed: %v", err)
	}

	if claims, ok := token.Claims.(*Token); ok && token.Valid {
		if claims.UserID != userID {
			t.Errorf("Expected UserID to be %d, but got %d", userID, claims.UserID)
		}
	} else {
		t.Errorf("Token is not valid")
	}
}

func TestGetUserFromContext(t *testing.T) {
	// Given
	userID := uint(1)
	ctx := context.WithValue(context.Background(), ContextUserKey, userID)

	// When
	actualUserID, ok := GetUserFromContext(ctx)

	// Then
	if !ok {
		t.Errorf("User was not found in context")
	}

	if actualUserID != userID {
		t.Errorf("Expected UserID to be %d, but got %d", userID, actualUserID)
	}
}
