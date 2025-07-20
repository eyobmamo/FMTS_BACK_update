package initiator

import (
	constructor "FMTS/internal/user/adapter/outbound/persistance"
	vihicle_persistance "FMTS/internal/vehicle/adapter/outbound/persistance"

	"FMTS/internal/user/port/outbound"
	vihicle_port "FMTS/internal/vehicle/port/outbound"
	"FMTS/pkg/utils"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Persistence struct {
	UserPersistence    outbound.UserRepoOutboundPort
	VehivlePersistence vihicle_port.VehicleRepo
}

func InitPersistence(client *mongo.Client, DB_name string, logger utils.Logger) Persistence {
	collectionNames := []string{
		"users",
		"vehicle",
	}

	return Persistence{
		UserPersistence:    constructor.InitUserRepo(client, DB_name, collectionNames[0], logger),
		VehivlePersistence: vihicle_persistance.InitVehicleRepo(client, DB_name, collectionNames[1], logger),
	}
}
