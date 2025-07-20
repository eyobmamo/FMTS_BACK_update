package initiator

import (
	userApplication "FMTS/internal/user/application"
	"FMTS/pkg/utils"

	vehicle_application "FMTS/internal/vehicle/application"
)

type Application struct {
	UserApp    userApplication.UserService
	VehicleApp vehicle_application.VehicleService
}

func InitApplication(domain Domain, logger utils.Logger) Application {
	return Application{
		UserApp:    userApplication.NewUserService(domain.UserDomain, logger),
		VehicleApp: vehicle_application.NewVehicleService(domain.VehicleDomain, logger),
	}

}
