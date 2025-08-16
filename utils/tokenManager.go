package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	entity "FMTS/internal/user/domain/entity"

	"github.com/golang-jwt/jwt/v5"
)

// UserPayload represents the user data that will be encrypted in the token
type UserPayload struct {
	UserID         string `json:"user_id"`
	UserRole       string `json:"user_role"`
	FullName       string `json:"full_name"`
	UserName       string `json:"username"`
	OrganizationID string `json:"organization_id,omitempty"`
	UserType       string `json:"user_type"` // "individual" or "company"
	PhoneNumber    string `json:"phone_number"`
}

type JWTManager interface {
	GenerateAccessToken(user *entity.User) (string, error)
	GenerateRefreshToken(user *entity.User) (string, error)
	VerifyAccessToken(token string) (string, error)  // returns userID
	VerifyRefreshToken(token string) (string, error) // returns userID
	AccessTokenTTL() int64
}

type jwtManager struct {
	secretKey       string
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
	encryptionKey   string
	encryptionIV    string
}

func NewJWTManager(secretKey string, accessTokenTTL, refreshTokenTTL time.Duration, encryptionKey, encryptionIV string) JWTManager {
	return &jwtManager{
		secretKey:       secretKey,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
		encryptionKey:   encryptionKey,
		encryptionIV:    encryptionIV,
	}
}

// encryptUserData encrypts user data using AES-256-CBC
func (j *jwtManager) encryptUserData(user *entity.User) (string, error) {
	// Determine user role based on customer type
	var userRole string
	switch user.CustomerType {
	case "individual":
		userRole = "USER"
	case "company":
		userRole = "FLEET_MANAGER" // Companies can manage fleets
	default:
		userRole = "USER" // Default to USER for safety
	}

	// Log the role assignment for debugging
	fmt.Printf("User %s has CustomerType: %s, assigned role: %s\n", user.ID.Hex(), user.CustomerType, userRole)

	// Create user payload
	payload := UserPayload{
		UserID:         user.ID.Hex(),
		UserRole:       userRole,
		FullName:       user.FullName,
		UserName:       user.FullName, // Using FullName as UserName
		OrganizationID: "",            // Set empty string if not available
		UserType:       string(user.CustomerType), // Using CustomerType as UserType
		PhoneNumber:    user.PhoneNumber,
	}

	// Convert to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	// Pad data to AES block size
	paddedData := pkcs7Pad(jsonData, aes.BlockSize)

	// Create cipher
	block, err := aes.NewCipher([]byte(j.encryptionKey))
	if err != nil {
		return "", err
	}

	// Encrypt
	ciphertext := make([]byte, len(paddedData))
	cipher.NewCBCEncrypter(block, []byte(j.encryptionIV)).CryptBlocks(ciphertext, paddedData)

	// Return hex encoded
	return hex.EncodeToString(ciphertext), nil
}

func (j *jwtManager) GenerateAccessToken(user *entity.User) (string, error) {
	// Encrypt user data
	encryptedData, err := j.encryptUserData(user)
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"data": encryptedData,
		"exp":  time.Now().Add(j.accessTokenTTL).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

func (j *jwtManager) GenerateRefreshToken(user *entity.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID.Hex(),
		"exp":     time.Now().Add(j.refreshTokenTTL).Unix(),
		"type":    "refresh",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

func (j *jwtManager) VerifyAccessToken(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secretKey), nil
	})
	if err != nil || !token.Valid {
		return "", errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid claims")
	}

	// Check if data field exists (for encrypted user data)
	if _, hasData := claims["data"]; hasData {
		// This is a new format token with encrypted data
		// The middleware will handle decryption and validation
		return "", errors.New("token contains encrypted data - use middleware for validation")
	}

	// Fallback to old format for backward compatibility
	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", errors.New("user_id missing in token")
	}
	return userID, nil
}

func (j *jwtManager) VerifyRefreshToken(tokenStr string) (string, error) {
	return j.verifyToken(tokenStr, true)
}

func (j *jwtManager) verifyToken(tokenStr string, isRefresh bool) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secretKey), nil
	})
	if err != nil || !token.Valid {
		return "", errors.New("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid claims")
	}
	if isRefresh {
		if t, ok := claims["type"].(string); !ok || t != "refresh" {
			return "", errors.New("not a refresh token")
		}
	}
	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", errors.New("user_id missing in token")
	}
	return userID, nil
}

func (j *jwtManager) AccessTokenTTL() int64 {
	return int64(j.accessTokenTTL.Seconds())
}

// pkcs7Pad adds PKCS7 padding to data
func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := make([]byte, padding)
	for i := range padtext {
		padtext[i] = byte(padding)
	}
	return append(data, padtext...)
}
