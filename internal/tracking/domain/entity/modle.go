package models

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type VehicleLocation struct {
	OwnerID   string    `json:"owner_id" bson:"owner_id"`
	VehicleID string    `json:"vehicle_id" bson:"vehicle_id"`
	Latitude  float64   `json:"latitude" bson:"latitude"`
	Longitude float64   `json:"longitude" bson:"longitude"`
	Speed     float64   `json:"speed,omitempty" bson:"speed,omitempty"`
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`
}

// Ozzo validation for VehicleLocation
func (v VehicleLocation) Validate() error {
	return validation.ValidateStruct(&v,
		validation.Field(&v.OwnerID,
			validation.Required,
			is.Alphanumeric,
		),
		validation.Field(&v.VehicleID,
			validation.Required,
			is.Alphanumeric,
		),
		validation.Field(&v.Latitude,
			validation.Required,
			validation.Min(-90.0),
			validation.Max(90.0),
		),
		validation.Field(&v.Longitude,
			validation.Required,
			validation.Min(-180.0),
			validation.Max(180.0),
		),
		validation.Field(&v.Speed,
			validation.Min(0.0), // optional, but if present must be >= 0
		),
		validation.Field(&v.Timestamp,
			validation.Required,
		),
	)
}

type VehicleID struct {
	VehicleID string `json:"vehicle_id" bson:"vehicle_id"`
}

// Ozzo validation for VehicleID
func (v VehicleID) Validate() error {
	return validation.ValidateStruct(&v,
		validation.Field(&v.VehicleID,
			validation.Required,
			is.Alphanumeric,
		),
	)
}
