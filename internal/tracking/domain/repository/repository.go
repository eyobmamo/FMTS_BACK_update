package service

import (
	entity "FMTS/internal/tracking/domain/entity"
	// "FMTS/utils"
	"context"
)

type DomainTracker interface {
	UpdateLocation(ctx context.Context, location entity.VehicleLocation) (entity.VehicleLocation, error)
	GetLatestVehicleLocationByID(ctx context.Context, vehicleID string) (entity.VehicleLocation, error)
	GetLatestVehicleLocationsByUserID(ctx context.Context, userID string) ([]*entity.VehicleLocation, error)
}
