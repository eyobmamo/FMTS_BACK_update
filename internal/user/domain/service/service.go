package service

import (
	"errors"
	"time"

	model "FMTS/internal/user/domain/entity"
	"FMTS/internal/user/domain/repository"
	"FMTS/pkg/utils"
	// "go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/v2/bson"
	// "go.mongodb.org/mongo-driver/v2/mongo"
)

type UserDomain struct {
	userRepo repository.UserRepo
	logger   utils.Logger
}

func NewUserDomainService(repo repository.UserRepo, logger utils.Logger) UserService {
	return &UserDomain{
		userRepo: repo,
		logger:   logger,
	}
}

type UserService interface {
	FindByEmailOrPhone(email string, phone string) (bool, error)
	CreateUser(user model.User) (*model.User, error)
	FindByID(id string) (*model.User, error)
	FindAll() ([]*model.User, error)
	UpdateUser(user model.User) error
	UpdateDelete(user *model.User, id string) error
}

// Check for existing user by email or phone
func (u *UserDomain) FindByEmailOrPhone(email string, phone string) (bool, error) {
	existingUser, err := u.userRepo.FindByEmailOrPhone(email, phone)
	if err != nil {
		u.logger.Errorf("[FindByEmailOrPhone] DB error: %v", err)
		return false, err
	}
	return existingUser != nil, nil
}

// Create a new user
func (u *UserDomain) CreateUser(user model.User) (*model.User, error) {
	user.ID = bson.NewObjectID()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.IsDeleted = false
	user.IsDisabled = false
	user.IsVerified = false

	createdUser, err := u.userRepo.CreateUser(user)
	if err != nil {
		u.logger.Errorf("[CreateUser] failed to create user: %v", err)
		return nil, err
	}
	return createdUser, nil
}

// Find user by ID
func (u *UserDomain) FindByID(id string) (*model.User, error) {
	user, err := u.userRepo.FindByID(id)
	if err != nil {
		u.logger.Errorf("[FindByID] error: %v", err)
		return nil, err
	}
	if user == nil || user.IsDeleted {
		return nil, errors.New("user not found or deleted")
	}
	return user, nil
}

// List all users
func (u *UserDomain) FindAll() ([]*model.User, error) {
	users, err := u.userRepo.FindAllUser()
	if err != nil {
		u.logger.Errorf("[FindAll] error: %v", err)
		return nil, err
	}
	return users, nil
}

// Update existing user
func (u *UserDomain) UpdateUser(user model.User) error {
	user.UpdatedAt = time.Now()

	if err := u.userRepo.UpdateUser(user); err != nil {
		u.logger.Errorf("[UpdateUser] error updating user: %v", err)
		return err
	}
	return nil
}

// Soft delete user
func (u *UserDomain) UpdateDelete(user *model.User, id string) error {
	user, err := u.userRepo.FindByID(id)
	if err != nil {
		u.logger.Errorf("[UpdateDelete] find error: %v", err)
		return err
	}

	if user.IsDeleted {
		return errors.New("user already deleted")
	}

	user.IsDeleted = true
	user.UpdatedAt = time.Now()

	if err := u.userRepo.UpdateSoftDelete(id); err != nil {
		u.logger.Errorf("[UpdateDelete] update error: %v", err)
		return err
	}
	return nil
}
