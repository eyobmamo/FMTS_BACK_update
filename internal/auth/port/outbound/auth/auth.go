package outbound

// import (
// 	entity "FMTS/internal/user/domain/entity"
// )

// UserRepo handles user persistence.
// type UserRepo interface {
// 	FindByEmail(email string) (*entity.User, error)
// 	FindByID(id string) (*entity.User, error)
// 	CreateUser(user entity.User) (*entity.User, error)
// }

// TokenRepo handles refresh token persistence and validation.
type TokenRepo interface {
	StoreRefreshToken(userID string, refreshToken string) error
	ValidateRefreshToken(userID string, refreshToken string) (bool, error)
	DeleteRefreshToken(userID string, refreshToken string) error
	InvalidateRefreshToken(refreshToken string) error
}
