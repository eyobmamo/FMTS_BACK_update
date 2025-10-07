package user

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type CreateUserRequest struct {
	FullName     string `json:"full_name" validate:"required"`
	FaydaID      string `json:"fayda_id" validate:"required"`
	Email        string `json:"email" validate:"required,email"`
	PhoneNumber  string `json:"phone_number" validate:"required,e164"`
	CustomerType string `json:"customer_type" validate:"required,oneof=individual company"`
}

// type UpdateUserRequest struct {
// 	FullName     *string `json:"full_name" validate:"required"`
// 	FaydaID      *string `json:"fayda_id" validate:"required"`
// 	Email        *string `json:"email" validate:"required,email"`
// 	PhoneNumber  *string `json:"phone_number" validate:"required,e164"`
// 	CustomerType *string `json:"customer_type" validate:"required,oneof=individual company"`
// }

func (c *CreateUserRequest) Validate() error {
	return validation.ValidateStruct(
		c,
		validation.Field(&c.FullName, validation.Required),
		validation.Field(&c.FaydaID, validation.Required),
		validation.Field(&c.Email, validation.Required, is.Email),
		validation.Field(&c.PhoneNumber, validation.Required, is.E164),
		validation.Field(&c.CustomerType, validation.Required, validation.In("individual", "company", "admin")),
	)
}

type UpdateUserRequest struct {
	FullName     *string `json:"full_name"`
	FaydaID      *string `json:"fayda_id"`
	Email        *string `json:"email"`
	PhoneNumber  *string `json:"phone_number"`
	CustomerType *string `json:"customer_type"`
}

func (u *UpdateUserRequest) Validate() error {
	var fullName, faydaID, email, phoneNumber, customerType string
	if u.FullName != nil {
		fullName = *u.FullName
	}
	if u.FaydaID != nil {
		faydaID = *u.FaydaID
	}
	if u.Email != nil {
		email = *u.Email
	}
	if u.PhoneNumber != nil {
		phoneNumber = *u.PhoneNumber
	}
	if u.CustomerType != nil {
		customerType = *u.CustomerType
	}

	return validation.ValidateStruct(
		struct {
			FullName     string `validate:"omitempty,min=2,max=100"`
			FaydaID      string `validate:"omitempty,min=5,max=50"`
			Email        string `validate:"omitempty,email"`
			PhoneNumber  string `validate:"omitempty,e164"`
			CustomerType string `validate:"omitempty,oneof=individual company"`
		}{
			FullName:     fullName,
			FaydaID:      faydaID,
			Email:        email,
			PhoneNumber:  phoneNumber,
			CustomerType: customerType,
		},
	)
}
