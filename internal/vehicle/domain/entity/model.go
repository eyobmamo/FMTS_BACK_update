package models

import (
	// "errors"
	"time"
)

type Vehicle struct {
	ID               string      `bson:"_id,omitempty" json:"id"`
	OwnerID          string      `bson:"owner_id" json:"owner_id"`
	OwnerType        OwnerType   `bson:"owner_type" json:"owner_type"`
	PlateNumber      string      `bson:"plate_number" json:"plate_number"`
	VehicleType      VehicleType `bson:"vehicle_type" json:"vehicle_type"`
	Model            string      `bson:"model" json:"model"`
	Manufacturer     string      `bson:"manufacturer,omitempty" json:"manufacturer,omitempty"`
	Year             int         `bson:"year,omitempty" json:"year,omitempty"`
	Color            string      `bson:"color,omitempty" json:"color,omitempty"`
	DriverName       string      `bson:"driver_name,omitempty" json:"driver_name,omitempty"`
	DriverPhone      string      `bson:"driver_phone,omitempty" json:"driver_phone,omitempty"`
	ImageURL         string      `bson:"image_url,omitempty" json:"image_url,omitempty"`
	CurrentlyTracked bool        `bson:"currently_tracked" json:"currently_tracked"`
	IsDeleted        bool        `bson:"is_deleted" json:"is_deleted"`
	IsDisabled       bool        `bson:"is_disabled" json:"is_disabled"`
	DisabledReason   string      `bson:"disabled_reason,omitempty" json:"disabled_reason,omitempty"`
	NotTrackedReason string      `bson:"not_tracked_reason,omitempty" json:"not_tracked_reason,omitempty"`
	CreatedAt        time.Time   `bson:"created_at" json:"created_at"`
	UpdatedAt        time.Time   `bson:"updated_at" json:"updated_at"`
}

// OwnerType custom string type with predefined values
type OwnerType string

const (
	OwnerTypePrivate    OwnerType = "private"
	OwnerTypePrivateOrg OwnerType = "private_organization"
	OwnerTypeGovernment OwnerType = "government"
	OwnerTypeNGO        OwnerType = "ngo"
)

// func (o OwnerType) IsValid() error {
// 	switch o {
// 	case OwnerTypePrivate, OwnerTypePrivateOrg, OwnerTypeGovernment, OwnerTypeNGO:
// 		return nil
// 	}
// 	return errors.New("invalid OwnerType")
// }

// VehicleType custom string type with predefined values
type VehicleType string

const (
	VehicleTypeSedan      VehicleType = "sedan"
	VehicleTypeSUV        VehicleType = "suv"
	VehicleTypeTruck      VehicleType = "truck"
	VehicleTypeVan        VehicleType = "van"
	VehicleTypeMotorcycle VehicleType = "motorcycle"
	VehicleTypeBus        VehicleType = "bus"
)

// func (v VehicleType) IsValid() error {
// 	switch v {
// 	case VehicleTypeSedan, VehicleTypeSUV, VehicleTypeTruck, VehicleTypeVan, VehicleTypeMotorcycle, VehicleTypeBus:
// 		return nil
// 	}
// 	return errors.New("invalid VehicleType")
// }
