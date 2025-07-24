package tracker

import (
	"net/http"

	inbound "FMTS/internal/tracking/port/inbound"
	route "FMTS/internal/user/adapter"
	"FMTS/internal/user/application/middleware"
	// role "FMTS/pkg/utils"

	"github.com/go-chi/chi/v5"
)

func InitUserRoutes(router chi.Router, userHandler inbound.TrackerPortHandler, authMiddleware middleware.AuthMiddleware) {
	router.Route("/tracker", func(r chi.Router) {
		routes := []route.Route{
			{
				Method:  http.MethodPost,
				Path:    "/",
				Handler: userHandler.UpdateLocation,
				Middlewares: []func(http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{"ADMIN"}),
				},
			},
			// {
			// 	Method:  http.MethodGet,
			// 	Path:    "/{id}",
			// 	Handler: userHandler.GetUserByID,
			// 	Middlewares: []func(http.Handler) http.Handler{
			// 		authMiddleware.AuthenticateToken,
			// 	},
			// },
			// {
			// 	Method:  http.MethodGet,
			// 	Path:    "/",
			// 	Handler: userHandler.ListUsers,
			// 	Middlewares: []func(http.Handler) http.Handler{
			// 		authMiddleware.AuthenticateToken,
			// 		authMiddleware.AccessControl([]string{"ADMIN"}),
			// 	},
			// },
			// {
			// 	Method:  http.MethodPatch,
			// 	Path:    "/{id}",
			// 	Handler: userHandler.UpdateUser,
			// 	Middlewares: []func(http.Handler) http.Handler{
			// 		authMiddleware.AuthenticateToken,
			// 	},
			// },
			// {
			// 	Method:  http.MethodDelete,
			// 	Path:    "/{id}",
			// 	Handler: userHandler.DeleteUser,
			// 	Middlewares: []func(http.Handler) http.Handler{
			// 		authMiddleware.AuthenticateToken,
			// 		authMiddleware.AccessControl([]string{"ADMIN"}),
			// 	},
			// },
		}

		route.RegisterRoutes(r, routes)
	})
}
