package initiator

import (
	userService "FMTS/internal/user/domain/service"
	vehicle_service "FMTS/internal/vehicle/domain/service"

	"FMTS/pkg/utils"
)

type Domain struct {
	UserDomain    userService.UserService
	VehicleDomain vehicle_service.VehicleService
}

func InitDomain(persistence Persistence, logger utils.Logger) Domain {
	return Domain{
		UserDomain:    userService.NewUserDomainService(persistence.UserPersistence, logger),
		VehicleDomain: vehicle_service.NewVehicleDomainService(persistence.VehivlePersistence, logger),
	}
}
