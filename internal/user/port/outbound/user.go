package outbound

import (
	model "FMTS/internal/user/domain/entity"
)

type UserRepoOutboundPort interface {
	FindByEmailOrPhone(email string, phone string) (*model.User, error)
	CreateUser(user model.User) (*model.User, error)
	FindByID(id string) (*model.User, error)
	FindAllUser() ([]*model.User, error)
	UpdateUser(model.User) error
	UpdateSoftDelete(id string) error
}
