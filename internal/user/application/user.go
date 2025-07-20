package user

import (
	"errors"
	"time"

	model "FMTS/internal/user/domain/entity"
	domain "FMTS/internal/user/domain/service"
	"FMTS/pkg/utils"
	// "go.mongodb.org/mongo-driver/bson/primitive"
)

// UserService defines all the business use cases for the User.
type UserService interface {
	CreateUser(req CreateUserRequest, createdBy string) (*model.User, error)
	GetUserByID(id string) (*model.User, error)
	ListUsers() ([]*model.User, error)
	UpdateUser(id string, req UpdateUserRequest) (*model.User, error)
	DeleteUser(id string) error
}

type userServiceImpl struct {
	domain domain.UserService
	logger utils.Logger
}

// Constructor
func NewUserService(domain domain.UserService, logger utils.Logger) UserService {
	return &userServiceImpl{
		domain: domain,
		logger: logger,
	}
}

// CreateUser handles registration logic with validation and persistence
func (s *userServiceImpl) CreateUser(req CreateUserRequest, createdBy string) (*model.User, error) {
	// Step 1: Validate input
	if err := req.Validate(); err != nil {
		s.logger.Warnf("[CreateUser] validation failed: %v", err)
		return nil, err
	}

	// Step 2: Check for duplicates
	existing, _ := s.domain.FindByEmailOrPhone(req.Email, req.PhoneNumber)
	if existing != false {
		return nil, errors.New("user already exists with provided email or phone")
	}

	// Step 3: Build user entity
	user := model.User{
		// ID:           bson.NewObjectID(),
		FullName:     req.FullName,
		FaydaID:      req.FaydaID,
		Email:        req.Email,
		PhoneNumber:  req.PhoneNumber,
		CustomerType: model.CustomerTypeIndividual,
		IsVerified:   false,
		IsDisabled:   false,
		IsDeleted:    false,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Step 4: Save to database
	CreatedUser, err := s.domain.CreateUser(user)
	if err != nil {
		s.logger.Errorf("[CreateUser] failed to save user: %v", err)
		return nil, err
	}

	return CreatedUser, nil
}

// GetUserByID returns a user by ObjectID
func (s *userServiceImpl) GetUserByID(id string) (*model.User, error) {
	user, err := s.domain.FindByID(id)
	if err != nil {
		s.logger.Errorf("[GetUserByID] error: %v", err)
		return nil, err
	}
	return user, nil
}

// ListUsers fetches all users
func (s *userServiceImpl) ListUsers() ([]*model.User, error) {
	users, err := s.domain.FindAll()
	if err != nil {
		s.logger.Errorf("[ListUsers] error: %v", err)
		return nil, err
	}
	return users, nil
}

// UpdateUser updates allowed fields
func (s *userServiceImpl) UpdateUser(id string, req UpdateUserRequest) (*model.User, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	user, err := s.domain.FindByID(id)
	if err != nil {
		return nil, err
	}
	if req.FullName != nil {
		user.FullName = *req.FullName
	}

	if req.PhoneNumber != nil {
		user.PhoneNumber = *req.PhoneNumber
	}
	if req.FullName != nil && *req.FullName != "" {
		user.FullName = *req.FullName
	}
	if req.FaydaID != nil && *req.FaydaID != "" {
		user.FaydaID = *req.FaydaID
	}
	if req.Email != nil && *req.Email != "" {
		user.Email = *req.Email
	}
	if req.PhoneNumber != nil && *req.PhoneNumber != "" {
		user.PhoneNumber = *req.PhoneNumber
	}
	// if req.CustomerType != nil && *req.CustomerType != "" {
	// 	user.CustomerType = *req.CustomerType
	// }
	user.UpdatedAt = time.Now()

	if err := s.domain.UpdateUser(*user); err != nil {
		return nil, err
	}

	return user, nil
}

// DeleteUser marks user as deleted
func (s *userServiceImpl) DeleteUser(id string) error {
	user, err := s.domain.FindByID(id)
	if err != nil {
		return err
	}

	user.IsDeleted = true
	user.UpdatedAt = time.Now()

	return s.domain.UpdateDelete(user, id)
}
