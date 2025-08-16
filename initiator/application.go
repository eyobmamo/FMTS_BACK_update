package initiator

import (
	userApplication "FMTS/internal/user/application"
	"FMTS/pkg/utils"

	tracker_application "FMTS/internal/tracking/application"
	// "FMTS/internal/tracking/domain/service"
	userAuth_application "FMTS/internal/auth/application"
	vehicle_application "FMTS/internal/vehicle/application"
)

type Application struct {
	UserApp     userApplication.UserService
	VehicleApp  vehicle_application.VehicleService
	TrackerApp  tracker_application.TrackerApplication
	AuthUserApp userAuth_application.AuthService
}

func InitApplication(domain Domain, logger utils.Logger) Application {
	return Application{
		UserApp:     userApplication.NewUserService(domain.UserDomain, logger),
		VehicleApp:  vehicle_application.NewVehicleService(domain.VehicleDomain, logger),
		TrackerApp:  tracker_application.NewTrackerApplicationService(domain.TrackerDomain, logger),
		AuthUserApp: userAuth_application.NewAuthService(domain.AuthUserDomain, logger),
	}

}
