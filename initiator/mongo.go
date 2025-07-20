package initiator

import (
	"FMTS/pkg/utils"

	// "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

func InitMongo(logger utils.Logger) *mongo.Client {
	mongoClient, err := ConnectToMongoDB()
	if err != nil {
		logger.Fatalf("failed to connect to mongo %v", err)
	}

	return mongoClient
}
