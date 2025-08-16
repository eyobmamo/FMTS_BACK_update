package initiator

import (
	tracking_persistance "FMTS/internal/tracking/adapter/outbound/mongo"
	constructor "FMTS/internal/user/adapter/outbound/persistance"
	vihicle_persistance "FMTS/internal/vehicle/adapter/outbound/persistance"

	tracking_port "FMTS/internal/tracking/port/outbound"
	"FMTS/internal/user/port/outbound"

	auth_persistance "FMTS/internal/auth/adapter/outbound/persistance/user"
	token_repo "FMTS/internal/auth/adapter/outbound/persistance/token"
	token "FMTS/internal/auth/port/outbound/auth"
	auth "FMTS/internal/auth/port/outbound/user"

	vihicle_port "FMTS/internal/vehicle/port/outbound"

	"FMTS/pkg/utils"

	config "FMTS/config"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Persistence struct {
	UserPersistence     outbound.UserRepoOutboundPort
	VehivlePersistence  vihicle_port.VehicleRepo
	TrackingPersistence tracking_port.TimescaleTrackerRepo
	AuthPersistance     token.TokenRepo
	AuthUserPersistance auth.UserRepo
}

var DB_URL = config.LoadConfig()

func InitPersistence(client *mongo.Client, DB_name string, logger utils.Logger) Persistence {
	collectionNames := []string{
		"users",
		"vehicle",
		"tracking",
		"tokens",
	}

	return Persistence{
		UserPersistence:     constructor.InitUserRepo(client, DB_name, collectionNames[0], logger),
		VehivlePersistence:  vihicle_persistance.InitVehicleRepo(client, DB_name, collectionNames[1], logger),
		TrackingPersistence: tracking_persistance.NewTimescaleTrackerRepo(config.ConnectSupabasePool(DB_URL)),
		AuthPersistance:     token_repo.InitTokenRepo(client, DB_name, collectionNames[3], logger),
		AuthUserPersistance: auth_persistance.NewUserAuthRepo(client, DB_name, collectionNames[0], logger),
	}
}
