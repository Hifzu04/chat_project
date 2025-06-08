// File: utils/utils.go
//hash password

package utils

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes the given plaintext password using bcrypt.
//
// It returns the hashed password (as a string) or an error if hashing fails.
func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// CheckPasswordHash compares a plaintext password with a bcrypt-hashed password.
// It returns true if they match, false otherwise.
func CheckPasswordHash(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// JSONResponse is a helper to write a JSON response with a given status code.
// `payload` must be JSON-marshalable.
func JSONResponse(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	// encode the payload
	json.NewEncoder(w).Encode(payload)
}
