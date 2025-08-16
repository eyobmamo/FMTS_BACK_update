package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a plain password using bcrypt.
func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

// CheckPasswordHash compares a plain password with a hashed password.
// import "fmt"

// PrintPasswords prints the raw and hashed password for debugging.
func PrintPasswords(raw, hashed string) {
	fmt.Printf("Raw password: %s\n", raw)
	fmt.Printf("Hashed password: %s\n", hashed)
}
func CheckPasswordHash(password, hash string) bool {
	fmt.Printf("Raw password: %s\n", password)
	fmt.Printf("Hashed password: %s\n", hash)
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
