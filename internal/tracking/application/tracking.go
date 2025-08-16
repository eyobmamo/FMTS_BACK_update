package application

import (
	entity "FMTS/internal/tracking/domain/entity"
	domain "FMTS/internal/tracking/domain/service"

	"FMTS/pkg/utils"
	"context"
)

type TrackerApplication interface {
	UpdateLocation(ctx context.Context, location entity.VehicleLocation) (entity.VehicleLocation, error)
	GetLatestVehicleLocationByID(ctx context.Context, vehicleID string) (entity.VehicleLocation, error)
	GetLatestVehicleLocationsByUserID(ctx context.Context, userID string) ([]*entity.VehicleLocation, error)
}
type TrackerApplicaionService struct {
	TrackerDomain domain.DomainTracker
	Logger        utils.Logger
}

func NewTrackerApplicationService(trackerDomain domain.DomainTracker, logger utils.Logger) TrackerApplication {
	return &TrackerApplicaionService{
		TrackerDomain: trackerDomain,
		Logger:        logger,
	}
}
func (s *TrackerApplicaionService) UpdateLocation(ctx context.Context, location entity.VehicleLocation) (entity.VehicleLocation, error) {
	updatedLocation, err := s.TrackerDomain.UpdateLocation(ctx, location)
	if err != nil {
		s.Logger.Errorf("[UpdateLocation] failed: %v", err)
		return entity.VehicleLocation{}, err
	}
	return updatedLocation, nil
}

func (s *TrackerApplicaionService) GetLatestVehicleLocationByID(ctx context.Context, vehicleID string) (entity.VehicleLocation, error) {
	location, err := s.TrackerDomain.GetLatestVehicleLocationByID(ctx, vehicleID)
	if err != nil {
		s.Logger.Errorf("[GetLatestVehicleLocationByID] failed: %v", err)
		return entity.VehicleLocation{}, err
	}
	return location, nil
}

func (s *TrackerApplicaionService) GetLatestVehicleLocationsByUserID(ctx context.Context, userID string) ([]*entity.VehicleLocation, error) {
	locations, err := s.TrackerDomain.GetLatestVehicleLocationsByUserID(ctx, userID)
	if err != nil {
		s.Logger.Errorf("[GetLatestVehicleLocationsByUserID] failed: %v", err)
		return nil, err
	}
	return locations, nil
}
