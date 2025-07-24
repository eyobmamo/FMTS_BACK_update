package repository

import (
	model "FMTS/internal/vehicle/domain/entity"
)

// VehicleRepo abstracts database operations for the Vehicle entity
type VehicleRepo interface {
	FindByPlateNumber(plate string) (*model.Vehicle, error)
	CreateVehicle(vehicle model.Vehicle) (*model.Vehicle, error)
	FindByID(id string) (*model.Vehicle, error)
	FindAllVehicles(User_ID string) ([]*model.Vehicle, error)
	UpdateVehicle(vehicle model.Vehicle) error
	UpdateSoftDelete(id string) error
}
