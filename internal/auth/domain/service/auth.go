package service

import (
	"errors"
	"time"

	"FMTS/utils"

	application "FMTS/internal/auth/application/dto"
	token "FMTS/internal/auth/domain/repository/token"
	"FMTS/internal/auth/domain/repository/user"

	entity "FMTS/internal/user/domain/entity"
)

type AuthDomainService interface {
	FindByEmail(email string) (*entity.User, error)
	CreateUser(user *entity.User) (*entity.User, error)
	GenerateTokens(user *entity.User) (*application.AuthTokens, error)
	RefreshTokens(refreshToken string) (*application.AuthTokens, error)
	InvalidateRefreshToken(refreshToken string) error
}

type AuthDomain struct {
	userRepo   repository.UserRepo
	tokenRepo  token.TokenRepo
	logger     utils.Logger
	jwtManager utils.JWTManager
}

func NewAuthDomainService(userRepo repository.UserRepo, tokenRepo token.TokenRepo, logger utils.Logger, jwtManager utils.JWTManager) AuthDomainService {
	return &AuthDomain{
		userRepo:   userRepo,
		tokenRepo:  tokenRepo,
		logger:     logger,
		jwtManager: jwtManager,
	}
}

func (a *AuthDomain) FindByEmail(email string) (*entity.User, error) {
	user, err := a.userRepo.FindByEmail(email)
	if err != nil {
		a.logger.Errorf("[FindByEmail] error: %v", err)
		return nil, err
	}
	return user, nil
}

func (a *AuthDomain) CreateUser(user *entity.User) (*entity.User, error) {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.IsDeleted = false
	user.IsDisabled = false
	user.IsVerified = false

	createdUser, err := a.userRepo.CreateUser(*user)
	if err != nil {
		a.logger.Errorf("[CreateUser] failed to create user: %v", err)
		return nil, err
	}
	return createdUser, nil
}

func (a *AuthDomain) GenerateTokens(user *entity.User) (*application.AuthTokens, error) {
	accessToken, err := a.jwtManager.GenerateAccessToken(user)
	if err != nil {
		a.logger.Errorf("[GenerateTokens] access token error: %v", err)
		return nil, err
	}
	refreshToken, err := a.jwtManager.GenerateRefreshToken(user)
	if err != nil {
		a.logger.Errorf("[GenerateTokens] refresh token error: %v", err)
		return nil, err
	}
	// Optionally, store refresh token in DB
	if err := a.tokenRepo.StoreRefreshToken(user.ID.Hex(), refreshToken); err != nil {
		a.logger.Errorf("[GenerateTokens] store refresh token error: %v", err)
		return nil, err
	}
	return &application.AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    a.jwtManager.AccessTokenTTL(),
	}, nil
}

func (a *AuthDomain) RefreshTokens(refreshToken string) (*application.AuthTokens, error) {
	userID, err := a.jwtManager.VerifyRefreshToken(refreshToken)
	if err != nil {
		a.logger.Errorf("[RefreshTokens] invalid refresh token: %v", err)
		return nil, errors.New("invalid refresh token")
	}
	// Optionally, check if refresh token is in DB
	valid, err := a.tokenRepo.ValidateRefreshToken(userID, refreshToken)
	if err != nil || !valid {
		a.logger.Errorf("[RefreshTokens] refresh token not valid: %v", err)
		return nil, errors.New("refresh token not valid")
	}
	user, err := a.userRepo.FindByID(userID)
	if err != nil || user == nil {
		a.logger.Errorf("[RefreshTokens] user not found: %v", err)
		return nil, errors.New("user not found")
	}
	return a.GenerateTokens(user)
}

func (a *AuthDomain) InvalidateRefreshToken(refreshToken string) error {
	userID, err := a.jwtManager.VerifyRefreshToken(refreshToken)
	if err != nil {
		a.logger.Errorf("[InvalidateRefreshToken] invalid refresh token: %v", err)
		return errors.New("invalid refresh token")
	}
	if err := a.tokenRepo.DeleteRefreshToken(userID, refreshToken); err != nil {
		a.logger.Errorf("[InvalidateRefreshToken] delete error: %v", err)
		return err
	}
	return nil
}
