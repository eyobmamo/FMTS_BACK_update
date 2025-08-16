package application

import (
	"errors"
	"fmt"
	"time"

	dto "FMTS/internal/auth/application/dto"
	service "FMTS/internal/auth/domain/service"
	entity "FMTS/internal/user/domain/entity"

	"FMTS/utils"
)

// AuthService defines all the business use cases for authentication.
type AuthService interface {
	RegisterPassword(req dto.RegisterRequest) (*entity.User, error)
	Login(req dto.LoginRequest) (*dto.AuthTokens, error)
	RefreshToken(req dto.RefreshRequest) (*dto.AuthTokens, error)
	Logout(req dto.LogoutRequest) error
}

type authServiceImpl struct {
	domain service.AuthDomainService
	logger utils.Logger
}

// Constructor
func NewAuthService(domain service.AuthDomainService, logger utils.Logger) AuthService {
	return &authServiceImpl{
		domain: domain,
		logger: logger,
	}
}

// RegisterUser handles user registration with password hashing and persistence.
func (s *authServiceImpl) RegisterPassword(req dto.RegisterRequest) (*entity.User, error) {
	if err := req.Validate(); err != nil {
		s.logger.Warnf("[RegisterUser] validation failed: %v", err)
		return nil, err
	}

	existing, _ := s.domain.FindByEmail(req.Email)
	if existing == nil {
		return nil, errors.New("user not register with provided email")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		s.logger.Errorf("[RegisterUser] password hash error: %v", err)
		return nil, err
	}

	user := &entity.User{
		Email:          req.Email,
		HashedPassword: hashedPassword,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	createdUser, err := s.domain.CreateUser(user)
	if err != nil {
		s.logger.Errorf("[RegisterUser] failed to save user: %v", err)
		return nil, err
	}

	return createdUser, nil
}

// Login authenticates user and issues tokens.
func (s *authServiceImpl) Login(req dto.LoginRequest) (*dto.AuthTokens, error) {
	if err := req.Validate(); err != nil {
		s.logger.Warnf("[Login] validation failed: %v", err)
		return nil, err
	}

	user, err := s.domain.FindByEmail(req.Email)
	if err != nil || user == nil {
		fmt.Println(" step 1")

		return nil, errors.New("invalid credentials")
	}

	if !utils.CheckPasswordHash(req.Password, user.HashedPassword) {
		fmt.Println(" step 2")
		return nil, errors.New("invalid credentials")
	}

	tokens, err := s.domain.GenerateTokens(user)
	if err != nil {
		fmt.Println("step 3")
		s.logger.Errorf("[Login] token generation error: %v", err)
		return nil, err
	}

	return tokens, nil
}

// RefreshToken rotates refresh token and issues new access token.
func (s *authServiceImpl) RefreshToken(req dto.RefreshRequest) (*dto.AuthTokens, error) {
	if err := req.Validate(); err != nil {
		s.logger.Warnf("[RefreshToken] validation failed: %v", err)
		return nil, err
	}

	tokens, err := s.domain.RefreshTokens(req.RefreshToken)
	if err != nil {
		s.logger.Errorf("[RefreshToken] refresh error: %v", err)
		return nil, err
	}

	return tokens, nil
}

// Logout invalidates the refresh token.
func (s *authServiceImpl) Logout(req dto.LogoutRequest) error {
	if err := req.Validate(); err != nil {
		s.logger.Warnf("[Logout] validation failed: %v", err)
		return err
	}

	if err := s.domain.InvalidateRefreshToken(req.RefreshToken); err != nil {
		s.logger.Errorf("[Logout] invalidate error: %v", err)
		return err
	}

	return nil
}
