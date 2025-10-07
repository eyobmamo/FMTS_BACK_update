package service

import (
	"errors"
	"time"

	model "FMTS/internal/vehicle/domain/entity"
	"FMTS/internal/vehicle/domain/repository"
	"FMTS/pkg/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type VehicleDomain struct {
	vehicleRepo repository.VehicleRepo
	logger      utils.Logger
}

func NewVehicleDomainService(repo repository.VehicleRepo, logger utils.Logger) VehicleService {
	return &VehicleDomain{
		vehicleRepo: repo,
		logger:      logger,
	}
}

type VehicleService interface {
	FindByPlateNumber(plate string) (*model.Vehicle, error)
	CreateVehicle(vehicle model.Vehicle) (*model.Vehicle, error)
	FindByID(id string) (*model.Vehicle, error)
	FindAll(User_ID string) ([]*model.Vehicle, error)
	UpdateVehicle(vehicle model.Vehicle) (model.Vehicle, error)
	UpdateSoftDelete(id string) error
}

// Check for existing vehicle by plate number
func (v *VehicleDomain) FindByPlateNumber(plate string) (*model.Vehicle, error) {
	existingVehicle, err := v.vehicleRepo.FindByPlateNumber(plate)
	if err != nil {
		v.logger.Errorf("[FindByPlateNumber] DB error: %v", err)
		return nil, err
	}
	return existingVehicle, nil
}

// Create a new vehicle
func (v *VehicleDomain) CreateVehicle(vehicle model.Vehicle) (*model.Vehicle, error) {
	vehicle.ID = primitive.NewObjectID().Hex()
	vehicle.CreatedAt = time.Now()
	vehicle.UpdatedAt = time.Now()
	vehicle.IsDeleted = false
	vehicle.IsDisabled = false
	vehicle.CurrentlyTracked = false

	createdVehicle, err := v.vehicleRepo.CreateVehicle(vehicle)
	if err != nil {
		v.logger.Errorf("[CreateVehicle] failed to create vehicle: %v", err)
		return nil, err
	}
	return createdVehicle, nil
}

// Find vehicle by ID
func (v *VehicleDomain) FindByID(id string) (*model.Vehicle, error) {
	vehicle, err := v.vehicleRepo.FindByID(id)
	if err != nil {
		v.logger.Errorf("[FindByID] error: %v", err)
		return nil, err
	}
	if vehicle == nil || vehicle.IsDeleted {
		return nil, errors.New("vehicle not found or deleted")
	}
	return vehicle, nil
}

// List all vehicles
func (v *VehicleDomain) FindAll(User_ID string) ([]*model.Vehicle, error) {
	vehicles, err := v.vehicleRepo.FindAllVehicles(User_ID)
	if err != nil {
		v.logger.Errorf("[FindAll] error: %v", err)
		return nil, err
	}
	return vehicles, nil
}

// Update existing vehicle
func (v *VehicleDomain) UpdateVehicle(vehicle model.Vehicle) (model.Vehicle, error) {
	vehicle.UpdatedAt = time.Now()

	updated, err := v.vehicleRepo.UpdateVehicle(vehicle)
	if err != nil {
		v.logger.Errorf("[UpdateVehicle] error updating vehicle: %v", err)
		return model.Vehicle{}, err
	}
	return updated, nil
}

// Soft delete vehicle
func (v *VehicleDomain) UpdateSoftDelete(id string) error {
	vehicle, err := v.vehicleRepo.FindByID(id)
	if err != nil {
		v.logger.Errorf("[UpdateSoftDelete] find error: %v", err)
		return err
	}

	if vehicle.IsDeleted {
		return errors.New("vehicle already deleted")
	}

	vehicle.IsDeleted = true
	vehicle.UpdatedAt = time.Now()

	if err := v.vehicleRepo.UpdateSoftDelete(id); err != nil {
		v.logger.Errorf("[UpdateSoftDelete] update error: %v", err)
		return err
	}
	return nil
}
