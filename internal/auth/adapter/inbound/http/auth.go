package auth

import (
	"encoding/json"
	"net/http"

	auth "FMTS/internal/auth/application"
	dto "FMTS/internal/auth/application/dto"
	port "FMTS/internal/auth/port/inbound"
	"FMTS/utils"
)

type AuthHandler struct {
	authService auth.AuthService
	logger      utils.Logger
}

func NewAuthHandler(service auth.AuthService, logger utils.Logger) port.AuthHandler {
	return &AuthHandler{
		authService: service,
		logger:      logger,
	}
}

func (h *AuthHandler) RegisterPassword(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Errorf("[RegisterUser] decode error: %v", err)
		utils.SendErrorResponse(w, "invalid request format", http.StatusBadRequest, nil)
		return
	}
	if err := req.Validate(); err != nil {
		h.logger.Warnf("[RegisterUser] validation failed: %v", err)
		utils.SendErrorResponse(w, err.Error(), http.StatusBadRequest, nil)
		return
	}
	user, err := h.authService.RegisterPassword(req)
	if err != nil {
		h.logger.Errorf("[RegisterUser] service error: %v", err)
		utils.SendErrorResponse(w, err.Error(), http.StatusInternalServerError, nil)
		return
	}
	utils.WriteSuccessResponse(w, user, "User registered successfully")
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Errorf("[Login] decode error: %v", err)
		utils.SendErrorResponse(w, "invalid request format", http.StatusBadRequest, nil)
		return
	}
	tokens, err := h.authService.Login(req)
	if err != nil {
		h.logger.Errorf("[Login] service error: %v", err)
		utils.SendErrorResponse(w, "invalid credentials", http.StatusUnauthorized, nil)
		return
	}
	utils.WriteSuccessResponse(w, tokens, "Login successful")
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req dto.RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Errorf("[RefreshToken] decode error: %v", err)
		utils.SendErrorResponse(w, "invalid request format", http.StatusBadRequest, nil)
		return
	}
	tokens, err := h.authService.RefreshToken(req)
	if err != nil {
		h.logger.Errorf("[RefreshToken] service error: %v", err)
		utils.SendErrorResponse(w, "invalid refresh token", http.StatusUnauthorized, nil)
		return
	}
	utils.WriteSuccessResponse(w, tokens, "Token refreshed successfully")
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req dto.LogoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Errorf("[Logout] decode error: %v", err)
		utils.SendErrorResponse(w, "invalid request format", http.StatusBadRequest, nil)
		return
	}
	if err := h.authService.Logout(req); err != nil {
		h.logger.Errorf("[Logout] service error: %v", err)
		utils.SendErrorResponse(w, err.Error(), http.StatusInternalServerError, nil)
		return
	}
	utils.WriteSuccessResponse(w, nil, "Logout successful")
}
