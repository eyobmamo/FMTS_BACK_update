package middleware

import (
	"FMTS/pkg/utils"
	common "FMTS/utils"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	constant "FMTS/utils"
	"github.com/golang-jwt/jwt/v5"
)

type UserPayload struct {
	UserID         string `json:"user_id"`
	UserRole       string `json:"user_role"`
	FullName       string `json:"full_name"`
	UserName       string `json:"username"`
	OrganizationID string `json:"organization_id,omitempty"`
	UserType       string `json:"user_type"` // "individual" or "company"
	PhoneNumber    string `json:"phone_number"`
}

type AuthMiddleware interface {
	AccessControl(allowedRoles []string) func(http.Handler) http.Handler
	AuthenticateToken(next http.Handler) http.Handler
}

type authMiddleware struct {
	logger       utils.Logger
	JWTSecretKey string
	Key          string
	IV           string
}

func InitAuthMiddleware(secretKey, key, iv string, logger utils.Logger) AuthMiddleware {
	return &authMiddleware{
		JWTSecretKey: secretKey,
		Key:          key,
		IV:           iv,
		logger:       logger,
	}
}

// Role-based Access Control
func (a *authMiddleware) AccessControl(allowedRoles []string) func(http.Handler) http.Handler {
	roleSet := make(map[string]struct{}, len(allowedRoles))
	for _, role := range allowedRoles {
		roleSet[strings.ToUpper(role)] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, ok := r.Context().Value(constant.ContextKey("user_role")).(string)
			if !ok || role == "" {
				common.SendErrorResponse(w, "unauthorized: role missing", http.StatusUnauthorized, nil)
				return
			}
			if _, allowed := roleSet[strings.ToUpper(role)]; !allowed {
				common.SendErrorResponse(w, "access denied: invalid role", http.StatusForbidden, nil)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// Token Authentication
func (a *authMiddleware) AuthenticateToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		if tokenString == "" {
			common.SendErrorResponse(w, "missing access token", http.StatusUnauthorized, nil)
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("invalid signing method")
			}
			return []byte(a.JWTSecretKey), nil
		})

		if err != nil || !token.Valid {
			a.logger.Warnf("token invalid: %v", err)
			common.SendErrorResponse(w, "invalid or expired token", http.StatusUnauthorized, nil)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			common.SendErrorResponse(w, "invalid token claims", http.StatusUnauthorized, nil)
			return
		}

		encrypted, ok := claims["data"].(string)
		if !ok || encrypted == "" {
			common.SendErrorResponse(w, "token data missing", http.StatusUnauthorized, nil)
			return
		}

		decrypted, err := a.decryptUserData(encrypted)
		if err != nil {
			a.logger.Warnf("decrypt failed: %v", err)
			common.SendErrorResponse(w, "token decryption failed", http.StatusUnauthorized, nil)
			return
		}

		var user UserPayload
		if err := json.Unmarshal([]byte(decrypted), &user); err != nil {
			common.SendErrorResponse(w, "invalid token user data", http.StatusUnauthorized, nil)
			return
		}

		ctx := context.WithValue(r.Context(), constant.ContextKey("user_id"), user.UserID)
		ctx = context.WithValue(ctx, constant.ContextKey("user_role"), user.UserRole)
		ctx = context.WithValue(ctx, constant.ContextKey("user_type"), user.UserType)
		ctx = context.WithValue(ctx, constant.ContextKey("organization_id"), user.OrganizationID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

// AES-256 CBC Decryption
func (a *authMiddleware) decryptUserData(data string) (string, error) {
	key, iv := []byte(a.Key), []byte(a.IV)
	if len(key) != 32 || len(iv) != aes.BlockSize {
		return "", errors.New("invalid key or IV length")
	}

	ciphertext, err := hex.DecodeString(data)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	decrypted := make([]byte, len(ciphertext))
	cipher.NewCBCDecrypter(block, iv).CryptBlocks(decrypted, ciphertext)

	unpadded := pkcs7Unpad(decrypted, aes.BlockSize)
	if unpadded == nil {
		return "", errors.New("invalid padding on decrypted data")
	}
	return string(unpadded), nil
}

func pkcs7Unpad(data []byte, blockSize int) []byte {
	length := len(data)
	if length == 0 || length%blockSize != 0 {
		return nil
	}
	padding := int(data[length-1])
	return data[:length-padding]
}
