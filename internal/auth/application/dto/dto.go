package application

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type RegisterRequest struct {
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}

func (r RegisterRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Email, validation.Required, validation.Length(5, 50), is.Email),
		validation.Field(&r.Password, validation.Required, validation.Length(6, 100)),
	)
}

type LoginRequest struct {
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}

func (r LoginRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Email, validation.Required, validation.Length(5, 50), is.Email),
		validation.Field(&r.Password, validation.Required, validation.Length(6, 100)),
	)
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" bson:"refresh_token"`
}

func (r RefreshRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.RefreshToken, validation.Required, validation.Length(10, 500)),
	)
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" bson:"refresh_token"`
}

func (r LogoutRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.RefreshToken, validation.Required, validation.Length(10, 500)),
	)
}

type AuthTokens struct {
	AccessToken  string `json:"access_token" bson:"access_token"`
	RefreshToken string `json:"refresh_token" bson:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in" bson:"expires_in"` // seconds until expiration
}
