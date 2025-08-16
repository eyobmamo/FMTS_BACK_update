package initiator

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func ConnectToMongoDB() (*mongo.Client, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, relying on environment variables")
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		return nil, ErrMissingMongoURI
	}

	clientOpts := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(clientOpts)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	if err = client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return client, nil
}

var ErrMissingMongoURI = &MissingEnvVarError{"MONGO_URI"}

type MissingEnvVarError struct {
	VarName string
}

func (e *MissingEnvVarError) Error() string {
	return "missing environment variable: " + e.VarName
}
