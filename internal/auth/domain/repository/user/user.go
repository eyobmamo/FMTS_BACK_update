package repository

import (
	entity "FMTS/internal/user/domain/entity"
)

// UserRepo handles user persistence.
type UserRepo interface {
	FindByEmail(email string) (*entity.User, error)
	FindByID(id string) (*entity.User, error)
	CreateUser(user entity.User) (*entity.User, error)
}
