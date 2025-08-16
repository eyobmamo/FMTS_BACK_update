package initiator

import (
	"FMTS/pkg/utils"

	user_adapter "FMTS/internal/user/adapter/inbound/http"
	user_port "FMTS/internal/user/port/inbound/user"

	vehicle_adapter "FMTS/internal/vehicle/adapter/inbound/http"
	vehicle_port "FMTS/internal/vehicle/port/inbound"

	tracker_adapter "FMTS/internal/tracking/adapter/inbound/http"
	tracker_port "FMTS/internal/tracking/port/inbound"

	authUser_adapter "FMTS/internal/auth/adapter/inbound/http"
	authUser_port "FMTS/internal/auth/port/inbound"
)

type Adapter struct {
	UserAdapter     user_port.UserPortHandler
	VihicleAdapter  vehicle_port.VehiclePortInterface
	TrackerAdapter  tracker_port.TrackerPortHandler
	AuthUserAdapter authUser_port.AuthHandler
}

func InitAdapter(application Application, logger utils.Logger) Adapter {
	return Adapter{
		UserAdapter:     user_adapter.NewUserHandler(application.UserApp, logger),
		VihicleAdapter:  vehicle_adapter.NewVehicleHandler(application.VehicleApp, logger),
		TrackerAdapter:  tracker_adapter.NewTrackerHandler(nil, application.TrackerApp, logger),
		AuthUserAdapter: authUser_adapter.NewAuthHandler(application.AuthUserApp, logger),
	}
}
