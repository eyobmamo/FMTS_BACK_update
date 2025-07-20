package initiator

import (
	"FMTS/pkg/utils"

	user_adapter "FMTS/internal/user/adapter/inbound/http"
	user_port "FMTS/internal/user/port/inbound/user"

	vehicle_adapter "FMTS/internal/vehicle/adapter/inbound/http"
	vehicle_port "FMTS/internal/vehicle/port/inbound"
)

type Adapter struct {
	UserAdapter    user_port.UserPortHandler
	VihicleAdapter vehicle_port.VehiclePortInterface
}

func InitAdapter(application Application, logger utils.Logger) Adapter {
	return Adapter{
		UserAdapter:    user_adapter.NewUserHandler(application.UserApp, logger),
		VihicleAdapter: vehicle_adapter.NewVehicleHandler(application.VehicleApp, logger),
	}
}
