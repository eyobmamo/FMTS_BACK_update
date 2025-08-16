package initiator

import (
	// authToken_service "FMTS/internal/auth/domain/service"
	authUser_service "FMTS/internal/auth/domain/service"
	tracker_service "FMTS/internal/tracking/domain/service"
	userService "FMTS/internal/user/domain/service"
	vehicle_service "FMTS/internal/vehicle/domain/service"

	"FMTS/utils"
)

type Domain struct {
	UserDomain     userService.UserService
	VehicleDomain  vehicle_service.VehicleService
	TrackerDomain  tracker_service.DomainTracker
	AuthUserDomain authUser_service.AuthDomainService
	JWTRelated     utils.JWTManager
}

func InitDomain(persistence Persistence, logger utils.Logger, JWT utils.JWTManager) Domain {
	return Domain{
		UserDomain:     userService.NewUserDomainService(persistence.UserPersistence, logger),
		VehicleDomain:  vehicle_service.NewVehicleDomainService(persistence.VehivlePersistence, logger),
		TrackerDomain:  tracker_service.InitDomaintrakerservice(logger, persistence.TrackingPersistence),
		AuthUserDomain: authUser_service.NewAuthDomainService(persistence.AuthUserPersistance, persistence.AuthPersistance, logger, JWT),
	}
}
