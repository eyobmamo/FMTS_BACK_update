package initiator

import (
	"FMTS/internal/middleware"
	user_handler "FMTS/internal/user/adapter/inbound/http"
	vehicle_handler "FMTS/internal/vehicle/adapter/inbound/http"

	"FMTS/pkg/utils"

	"github.com/go-chi/chi/v5"
)

func InitRoutes(r chi.Router, adapter Adapter, secretKey, key, iv string, logger utils.Logger) {
	authMiddleware := middleware.InitAuthMiddleware(secretKey, key, iv, logger)

	r.Route("/api/v1/FMTS/", func(r chi.Router) {
		user_handler.InitUserRoutes(r, adapter.UserAdapter, authMiddleware)
		vehicle_handler.InitVehicleRoutes(r, adapter.VihicleAdapter, authMiddleware)

	})
}
