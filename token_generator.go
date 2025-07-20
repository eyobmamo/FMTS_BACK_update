package token

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserPayload struct {
	UserID         string `json:"user_id"`
	UserRole       string `json:"user_role"`
	FullName       string `json:"full_name"`
	UserName       string `json:"username"`
	OrganizationID string `json:"organization_id"`
	UserType       string `json:"user_type"` // "company" or "individual"
	PhoneNumber    string `json:"phone_number"`
}

// GenerateAdminToken generates JWT token with encrypted user payload
func GenerateAdminToken() (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET_KEY")
	aesKey := os.Getenv("AES_KEY")
	aesIV := os.Getenv("AES_IV")

	if jwtSecret == "" || aesKey == "" || aesIV == "" {
		return "", errors.New("JWT_SECRET_KEY, AES_KEY or AES_IV not set in env")
	}

	// 1. Create Admin UserPayload
	adminPayload := UserPayload{
		UserID:         "admin-001",
		UserRole:       "ADMIN",
		FullName:       "Admin User",
		UserName:       "admin",
		OrganizationID: "admin-org-001",
		UserType:       "company",
		PhoneNumber:    "+251900000000",
	}

	// 2. Marshal payload to JSON
	jsonData, err := json.Marshal(adminPayload)
	if err != nil {
		return "", fmt.Errorf("marshal error: %w", err)
	}

	// 3. Encrypt user payload
	encryptedData, err := encryptAES(string(jsonData), aesKey, aesIV)
	if err != nil {
		return "", fmt.Errorf("encryption error: %w", err)
	}

	// 4. Create JWT token with encrypted data
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"data": encryptedData,
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	})

	// 5. Sign token
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", fmt.Errorf("token sign error: %w", err)
	}

	return tokenString, nil
}

func encryptAES(plainText, keyString, ivString string) (string, error) {
	key, iv := []byte(keyString), []byte(ivString)
	if len(key) != 32 || len(iv) != aes.BlockSize {
		return "", errors.New("invalid AES key or IV length")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	plainBytes := pkcs7Pad([]byte(plainText), aes.BlockSize)
	cipherText := make([]byte, len(plainBytes))

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText, plainBytes)

	return hex.EncodeToString(cipherText), nil
}

func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	pad := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, pad...)
}
