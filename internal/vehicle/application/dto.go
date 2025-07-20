package vehicle

import (
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type CreateVehicleRequest struct {
	OwnerID      string      `json:"owner_id" validate:"required"`
	OwnerType    OwnerType   `json:"owner_type" validate:"required"`
	PlateNumber  string      `json:"plate_number" validate:"required"`
	VehicleType  VehicleType `json:"vehicle_type" validate:"required"`
	Model        string      `json:"model" validate:"required"`
	Manufacturer string      `json:"manufacturer,omitempty"`
	Year         int         `json:"year,omitempty"`
	Color        string      `json:"color,omitempty"`
	DriverName   string      `json:"driver_name,omitempty"`
	DriverPhone  string      `json:"driver_phone,omitempty"`
	ImageURL     string      `json:"image_url,omitempty"`
}

// Assume you have OwnerType and VehicleType as string aliases or custom types,
// and validation.In(...) will check for allowed values.

func (r CreateVehicleRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.OwnerID, validation.Required),
		// validation.Field(&r.OwnerType, validation.Required, validation.In("private_org", "private", "government", "ngo")), // Adjust string values to your enum
		validation.Field(&r.PlateNumber, validation.Required, validation.Length(1, 20)), // plate number length example
		// validation.Field(&r.VehicleType, validation.Required, validation.In("sedan", "truck", "van", "bus")),              // Adjust as needed
		validation.Field(&r.Model, validation.Required, validation.Length(1, 50)),

		// Optional fields, validate only if non-empty
		validation.Field(&r.Manufacturer, validation.NilOrNotEmpty, validation.Length(0, 50)),
		validation.Field(&r.Year, validation.NilOrNotEmpty, validation.Min(1886), validation.Max(2100)), // Cars invented ~1886, max arbitrary
		validation.Field(&r.Color, validation.NilOrNotEmpty, validation.Length(0, 30)),
		validation.Field(&r.DriverName, validation.NilOrNotEmpty, validation.Length(0, 100)),
		validation.Field(&r.DriverPhone, validation.NilOrNotEmpty, is.E164),
		validation.Field(&r.ImageURL, validation.NilOrNotEmpty, is.URL),
	)
}

type OwnerType string

const (
	OwnerTypePrivate    OwnerType = "private"
	OwnerTypePrivateOrg OwnerType = "private_organization"
	OwnerTypeGovernment OwnerType = "government"
	OwnerTypeNGO        OwnerType = "ngo"
)

type VehicleType string

const (
	VehicleTypeSedan      VehicleType = "sedan"
	VehicleTypeSUV        VehicleType = "suv"
	VehicleTypeTruck      VehicleType = "truck"
	VehicleTypeVan        VehicleType = "van"
	VehicleTypeMotorcycle VehicleType = "motorcycle"
	VehicleTypeBus        VehicleType = "bus"
)

type UpdateVehicleRequest struct {
	OwnerID          *string      `json:"owner_id,omitempty"`
	OwnerType        *OwnerType   `json:"owner_type,omitempty"`
	PlateNumber      *string      `json:"plate_number,omitempty"`
	VehicleType      *VehicleType `json:"vehicle_type,omitempty"`
	Model            *string      `json:"model,omitempty"`
	Manufacturer     *string      `json:"manufacturer,omitempty"`
	Year             *int         `json:"year,omitempty"`
	Color            *string      `json:"color,omitempty"`
	DriverName       *string      `json:"driver_name,omitempty"`
	DriverPhone      *string      `json:"driver_phone,omitempty"`
	ImageURL         *string      `json:"image_url,omitempty"`
	CurrentlyTracked *bool        `json:"currently_tracked,omitempty"`
	IsDisabled       *bool        `json:"is_disabled,omitempty"`
	DisabledReason   *string      `json:"disabled_reason,omitempty"`
	NotTrackedReason *string      `json:"not_tracked_reason,omitempty"`
}

func (r UpdateVehicleRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.OwnerID, validation.When(r.OwnerID != nil, validation.Length(1, 100))),
		validation.Field(&r.OwnerType, validation.When(r.OwnerType != nil, validation.In(
			"PrivateOrganization", "Private", "Government", "NGO")), // Adjust values based on your enum
		),
		validation.Field(&r.PlateNumber, validation.When(r.PlateNumber != nil, validation.Length(1, 20))),
		validation.Field(&r.VehicleType, validation.When(r.VehicleType != nil, validation.In(
			"Sedan", "SUV", "Truck", "Van", "Motorcycle")), // Adjust values based on your enum
		),
		validation.Field(&r.Model, validation.When(r.Model != nil, validation.Length(1, 50))),
		validation.Field(&r.Manufacturer, validation.When(r.Manufacturer != nil, validation.Length(1, 50))),
		validation.Field(&r.Year, validation.When(r.Year != nil, validation.By(func(value interface{}) error {
			year, ok := value.(*int)
			if !ok || year == nil {
				return nil // skip if nil
			}
			if *year < 1900 || *year > 2100 {
				return errors.New("year must be between 1900 and 2100")
			}
			return nil
		}))),
		validation.Field(&r.Color, validation.When(r.Color != nil, validation.Length(1, 30))),
		validation.Field(&r.DriverName, validation.When(r.DriverName != nil, validation.Length(1, 50))),
		validation.Field(&r.DriverPhone, validation.When(r.DriverPhone != nil, is.E164)),
		validation.Field(&r.ImageURL, validation.When(r.ImageURL != nil, is.URL)),
		validation.Field(&r.DisabledReason,
			validation.When(r.IsDisabled != nil && *r.IsDisabled, validation.Required.Error("disabled_reason is required when vehicle is disabled")),
		),

		// Optional: If NotTrackedReason is present, you can add some validation if needed
		// for example: max length or non-empty
		validation.Field(&r.NotTrackedReason,
			validation.When(r.NotTrackedReason != nil, validation.Length(1, 500)),
		))
}
