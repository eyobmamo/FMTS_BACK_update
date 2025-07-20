package vehicle_handler

import (
	"net/http"

	"FMTS/internal/user/application/middleware"
	route "FMTS/internal/vehicle/adapter"
	inbound "FMTS/internal/vehicle/port/inbound"

	"github.com/go-chi/chi/v5"
)

func InitVehicleRoutes(router chi.Router, vehicleHandler inbound.VehiclePortInterface, authMiddleware middleware.AuthMiddleware) {
	router.Route("/vehicles", func(r chi.Router) {
		routes := []route.Route{
			{
				Method:  http.MethodPost,
				Path:    "/",
				Handler: vehicleHandler.RegisterVehicle,
				Middlewares: []func(http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{"ADMIN", "FLEET_MANAGER"}),
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/{id}",
				Handler: vehicleHandler.GetVehicleByID,
				Middlewares: []func(http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/",
				Handler: vehicleHandler.ListVehicles,
				Middlewares: []func(http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{"ADMIN", "FLEET_MANAGER"}),
				},
			},
			{
				Method:  http.MethodPatch,
				Path:    "/{id}",
				Handler: vehicleHandler.UpdateVehicle,
				Middlewares: []func(http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{"ADMIN", "FLEET_MANAGER"}),
				},
			},
			{
				Method:  http.MethodDelete,
				Path:    "/{id}",
				Handler: vehicleHandler.DeleteVehicle,
				Middlewares: []func(http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{"ADMIN"}),
				},
			},
		}

		route.RegisterRoutes(r, routes)
	})
}
