package vehicle

import (
	"errors"
	"time"

	model "FMTS/internal/vehicle/domain/entity"
	domain "FMTS/internal/vehicle/domain/service"
	"FMTS/pkg/utils"
)

// VehicleService defines business use cases for Vehicle
type VehicleService interface {
	RegisterVehicle(req CreateVehicleRequest, createdBy string) (*model.Vehicle, error)
	GetVehicleByID(id string) (*model.Vehicle, error)
	ListVehicles() ([]*model.Vehicle, error)
	UpdateVehicle(id string, req UpdateVehicleRequest) (*model.Vehicle, error)
	DeleteVehicle(id string) error
}

type vehicleServiceImpl struct {
	domain domain.VehicleService
	logger utils.Logger
}

// Constructor
func NewVehicleService(domain domain.VehicleService, logger utils.Logger) VehicleService {
	return &vehicleServiceImpl{
		domain: domain,
		logger: logger,
	}
}

// RegisterVehicle handles validation, uniqueness check, and persists vehicle entity
func (s *vehicleServiceImpl) RegisterVehicle(req CreateVehicleRequest, createdBy string) (*model.Vehicle, error) {
	// 1. Validate input
	if err := req.Validate(); err != nil {
		s.logger.Warnf("[RegisterVehicle] validation failed: %v", err)
		return nil, err
	}

	// 2. Check for duplicate plate number
	existing, _ := s.domain.FindByPlateNumber(req.PlateNumber)
	if existing != nil {
		return nil, errors.New("vehicle already registered with this plate number")
	}

	// 3. Map request to entity
	vehicle := model.Vehicle{
		OwnerID:          req.OwnerID,
		OwnerType:        model.OwnerType(req.OwnerType),
		PlateNumber:      req.PlateNumber,
		VehicleType:      model.VehicleType(req.VehicleType),
		Model:            req.Model,
		Manufacturer:     req.Manufacturer,
		Year:             req.Year,
		Color:            req.Color,
		DriverName:       req.DriverName,
		DriverPhone:      req.DriverPhone,
		ImageURL:         req.ImageURL,
		CurrentlyTracked: false,
		IsDeleted:        false,
		IsDisabled:       false,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	// 4. Save to database via domain layer
	createdVehicle, err := s.domain.CreateVehicle(vehicle)
	if err != nil {
		s.logger.Errorf("[RegisterVehicle] failed to save vehicle: %v", err)
		return nil, err
	}

	return createdVehicle, nil
}

// GetVehicleByID fetches a vehicle entity by ID
func (s *vehicleServiceImpl) GetVehicleByID(id string) (*model.Vehicle, error) {
	vehicle, err := s.domain.FindByID(id)
	if err != nil {
		s.logger.Errorf("[GetVehicleByID] error: %v", err)
		return nil, err
	}
	return vehicle, nil
}

// ListVehicles returns all vehicles (pagination/filter can be added later)
func (s *vehicleServiceImpl) ListVehicles() ([]*model.Vehicle, error) {
	vehicles, err := s.domain.FindAll()
	if err != nil {
		s.logger.Errorf("[ListVehicles] error: %v", err)
		return nil, err
	}
	return vehicles, nil
}

// UpdateVehicle updates allowed fields for a vehicle
func (s *vehicleServiceImpl) UpdateVehicle(id string, req UpdateVehicleRequest) (*model.Vehicle, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	vehicle, err := s.domain.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields if present in the request
	if req.OwnerID != nil {
		vehicle.OwnerID = *req.OwnerID
	}
	// if req.OwnerType != nil {
	// 	vehicle.OwnerType = *req.OwnerType
	// }
	if req.PlateNumber != nil {
		vehicle.PlateNumber = *req.PlateNumber
	}
	// if req.VehicleType != nil {
	// 	vehicle.VehicleType = *req.VehicleType
	// }
	if req.Model != nil {
		vehicle.Model = *req.Model
	}
	if req.Manufacturer != nil {
		vehicle.Manufacturer = *req.Manufacturer
	}
	if req.Year != nil {
		vehicle.Year = *req.Year
	}
	if req.Color != nil {
		vehicle.Color = *req.Color
	}
	if req.DriverName != nil {
		vehicle.DriverName = *req.DriverName
	}
	if req.DriverPhone != nil {
		vehicle.DriverPhone = *req.DriverPhone
	}
	if req.ImageURL != nil {
		vehicle.ImageURL = *req.ImageURL
	}
	if req.CurrentlyTracked != nil {
		vehicle.CurrentlyTracked = *req.CurrentlyTracked
	}
	if req.IsDisabled != nil {
		vehicle.IsDisabled = *req.IsDisabled
	}
	if req.DisabledReason != nil {
		vehicle.DisabledReason = *req.DisabledReason
	}
	if req.NotTrackedReason != nil {
		vehicle.NotTrackedReason = *req.NotTrackedReason
	}

	vehicle.UpdatedAt = time.Now()

	if err := s.domain.UpdateVehicle(*vehicle); err != nil {
		return nil, err
	}

	return vehicle, nil
}

// DeleteVehicle marks vehicle as deleted (soft delete)
func (s *vehicleServiceImpl) DeleteVehicle(id string) error {
	vehicle, err := s.domain.FindByID(id)
	if err != nil {
		return err
	}

	vehicle.IsDeleted = true
	vehicle.UpdatedAt = time.Now()

	return s.domain.UpdateVehicle(*vehicle)
}
