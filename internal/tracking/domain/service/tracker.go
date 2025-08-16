package service

import (
	entity "FMTS/internal/tracking/domain/entity"
	repo "FMTS/internal/tracking/domain/repository"

	"FMTS/utils"
	"context"
)

type DomainTracker interface {
	UpdateLocation(ctx context.Context, location entity.VehicleLocation) (entity.VehicleLocation, error)
	GetLatestVehicleLocationByID(ctx context.Context, vehicleID string) (entity.VehicleLocation, error)
	GetLatestVehicleLocationsByUserID(ctx context.Context, UserID string) ([]*entity.VehicleLocation, error)
}

type DomainTrackerService struct {
	logger      utils.Logger
	trackerRepo repo.DomainTracker
}

func InitDomaintrakerservice(logger utils.Logger, trackerRepo repo.DomainTracker) *DomainTrackerService {
	return &DomainTrackerService{
		logger:      logger,
		trackerRepo: trackerRepo,
	}
}

func (s *DomainTrackerService) UpdateLocation(ctx context.Context, location entity.VehicleLocation) (entity.VehicleLocation, error) {
	return s.trackerRepo.UpdateLocation(ctx, location)
}

func (s *DomainTrackerService) GetLatestVehicleLocationByID(ctx context.Context, vehicleID string) (entity.VehicleLocation, error) {
	return s.trackerRepo.GetLatestVehicleLocationByID(ctx, vehicleID)
}

func (s *DomainTrackerService) GetLatestVehicleLocationsByUserID(ctx context.Context, userID string) ([]*entity.VehicleLocation, error) {
	return s.trackerRepo.GetLatestVehicleLocationsByUserID(ctx, userID)
}
