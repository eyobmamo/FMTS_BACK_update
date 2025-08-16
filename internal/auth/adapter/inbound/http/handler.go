package auth

import (
	"net/http"

	route "FMTS/internal/auth/adapter"
	inbound "FMTS/internal/auth/port/inbound"
	"FMTS/internal/user/application/middleware"
	// role "FMTS/pkg/utils"

	"github.com/go-chi/chi/v5"
)

func InitUserRoutes(router chi.Router, authHandler inbound.AuthHandler, authMiddleware middleware.AuthMiddleware) {
	router.Route("/auth", func(r chi.Router) {
		routes := []route.Route{
			{
				Method:      http.MethodPost,
				Path:        "/register",
				Handler:     authHandler.RegisterPassword,
				Middlewares: nil,
			},
			{
				Method:      http.MethodPost,
				Path:        "/login",
				Handler:     authHandler.Login,
				Middlewares: nil,
			},
			{
				Method:      http.MethodPost,
				Path:        "/Refresh",
				Handler:     authHandler.RefreshToken,
				Middlewares: nil,
			},
			{
				Method:  http.MethodPatch,
				Path:    "/logout",
				Handler: authHandler.Logout,
				Middlewares: []func(http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
				},
			},
		}

		route.RegisterRoutes(r, routes)
	})
}
