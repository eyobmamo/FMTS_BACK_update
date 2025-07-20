package model

import (
	// "go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	// "go.mongodb.org/mongo-driver/v2/mongo"
	"time"
)

type CustomerType string

const (
	CustomerTypeIndividual CustomerType = "individual"
	CustomerTypeCompany    CustomerType = "company"
)

type User struct {
	ID             bson.ObjectID `bson:"_id,omitempty" json:"id"`
	FullName       string        `bson:"full_name" json:"full_name" validate:"required"`
	FaydaID        string        `bson:"fayda_id" json:"fayda_id" validate:"required"`
	Email          string        `bson:"email" json:"email" validate:"required,email"`
	PhoneNumber    string        `bson:"phone_number" json:"phone_number" validate:"required,e164"` // E.164 format for phone
	CustomerType   CustomerType  `bson:"customer_type" json:"customer_type" validate:"required,oneof=individual company"`
	OTPExpiresAt   time.Time     `bson:"otp_expires_at" json:"otp_expires_at"`
	HashedPassword string        `bson:"hashed_password" json:"-" validate:"required"` // don’t expose in JSON
	IsVerified     bool          `bson:"is_verified" json:"is_verified"`
	IsDisabled     bool          `bson:"is_disabled" json:"is_disabled"`
	IsDeleted      bool          `bson:"is_deleted" json:"is_deleted"`
	CreatedAt      time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time     `bson:"updated_at" json:"updated_at"`
}
