package models

import (
	"context"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func TestGetUserFromContext(t *testing.T) {
	// Test case 1: Verify that GetUserFromContext() function returns the expected user ID and true boolean value.
	testID := uint(123)
	testContext := context.WithValue(context.Background(), ContextUserKey, testID)
	resultID, resultOK := GetUserFromContext(testContext)
	if resultID != testID {
		t.Errorf("Test case 1: expected GetUserFromContext() to return user ID %d, but got %d", testID, resultID)
	}
	if !resultOK {
		t.Errorf("Test case 1: expected GetUserFromContext() to return true, but got false")
	}

	// Test case 2: Verify that GetUserFromContext() function returns the expected zero user ID and false boolean value.
	testContext = context.Background()
	resultID, resultOK = GetUserFromContext(testContext)
	if resultID != 0 {
		t.Errorf("Test case 2: expected GetUserFromContext() to return user ID 0, but got %d", resultID)
	}
	if resultOK {
		t.Errorf("Test case 2: expected GetUserFromContext() to return false, but got true")
	}

}

func TestGetToken(t *testing.T) {
	// Set up test data
	testID := uint(123)
	jwtPassword := "my-secret-key"
	expTime := 60 // minutes
	InitJWTPassword(jwtPassword, expTime)

	// Test case 1: Verify that GetToken() function returns a non-empty token string and can be parsed.
	tokenString := GenerateToken(testID)
	if tokenString == "" {
		t.Errorf("Test case 1: expected GetToken() to return a non-empty token string, but got an empty string")
	}

	// Parse the token string
	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtPassword), nil
	})
	if err != nil {
		t.Errorf("Test case 1: expected parsed token error to be nil, but got %v", err)
	}
	if parsedToken == nil {
		t.Errorf("Test case 1: expected parsed token to be non-nil, but got nil")
	} else if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		if userID, exists := claims["user_id"].(uint); exists {
			t.Errorf("Test case 1: expected user ID in parsed token to be %d, but got %d", testID, userID)
		}
		if exp, exists := claims["exp"].(float64); exists {
			if expirationTime := time.Unix(int64(exp), 0); !expirationTime.After(time.Now()) {
				t.Errorf("Test case 1: expected expiration time in parsed token to be in the future, but got %v", expirationTime)
			}
		} else {
			t.Errorf("Test case 1: expected exp claim to exist in parsed token, but got nil")
		}
	} else {
		t.Errorf("Test case 1: expected parsed token claims to be of type MapClaims and valid, but got type %T and validity %v", parsedToken.Claims, parsedToken.Valid)
	}

}

func TestEncryptPassword(t *testing.T) {
	// Set up test data
	testPassword := "my-secret-password"

	// Test case 1: Verify that EncryptPassword() function returns a non-empty hashed password string.
	hashedPassword := EncryptPassword(testPassword)
	if hashedPassword == "" {
		t.Errorf("Test case 1: expected EncryptPassword() to return a non-empty hashed password string, but got an empty string")
	}

	// Test case 2: Verify that the hashed password can be used to verify the original password.
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))
	if err != nil {
		t.Errorf("Test case 2: expected password verification error to be nil, but got %v", err)
	}

}

func TestIsPasswordsEqual(t *testing.T) {
	// Set up test data
	testExistingPassword := "my-secret-password"
	testNewPassword := "my-new-password"

	// Generate a hashed password from the existing password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(testExistingPassword), bcrypt.DefaultCost)

	// Test case 1: Verify that IsPasswordsEqual() function returns true when the passwords match.
	result := IsPasswordsEqual(string(hashedPassword), testExistingPassword)
	if !result {
		t.Errorf("Test case 1: expected IsPasswordsEqual() to return true, but got false")
	}

	// Test case 2: Verify that IsPasswordsEqual() function returns false when the passwords don't match.
	result = IsPasswordsEqual(string(hashedPassword), testNewPassword)
	if result {
		t.Errorf("Test case 2: expected IsPasswordsEqual() to return false, but got true")
	}

}
