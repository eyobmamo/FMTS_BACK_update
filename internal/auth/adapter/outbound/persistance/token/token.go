package repo_token

import (
	"context"
	"errors"
	"time"

	"FMTS/internal/auth/domain/repository/token"
	dal "FMTS/internal/user/adapter/outbound/infra"
	"FMTS/pkg/utils"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// TokenModel represents a refresh token document.
type TokenModel struct {
	UserID       string    `bson:"user_id"`
	RefreshToken string    `bson:"refresh_token"`
	CreatedAt    time.Time `bson:"created_at"`
}

type TokenPersistence struct {
	tokenDal dal.MongoDal[TokenModel, TokenModel]
	logger   utils.Logger
}

var _ repository.TokenRepo = (*TokenPersistence)(nil)

func InitTokenRepo(client *mongo.Client, dbName string, collection string, logger utils.Logger) repository.TokenRepo {
	tokenDal := dal.NewMongoDal[TokenModel, TokenModel](client, dbName, collection)
	return &TokenPersistence{
		tokenDal: tokenDal,
		logger:   logger,
	}
}

func (t *TokenPersistence) StoreRefreshToken(userID string, refreshToken string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	token := TokenModel{
		UserID:       userID,
		RefreshToken: refreshToken,
		CreatedAt:    time.Now(),
	}
	_, err := t.tokenDal.InsertOne(ctx, token)
	if err != nil {
		t.logger.Errorf("[StoreRefreshToken] insert error: %v", err)
		return err
	}
	return nil
}

func (t *TokenPersistence) ValidateRefreshToken(userID string, refreshToken string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"user_id": userID, "refresh_token": refreshToken}
	projection := bson.M{}
	token, err := t.tokenDal.FindOne(ctx, filter, projection)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}
		t.logger.Errorf("[ValidateRefreshToken] find error: %v", err)
		return false, err
	}
	return token != nil, nil
}

func (t *TokenPersistence) DeleteRefreshToken(userID string, refreshToken string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"user_id": userID, "refresh_token": refreshToken}
	err := t.tokenDal.DeleteOne(ctx, filter)
	if err != nil {
		t.logger.Errorf("[DeleteRefreshToken] delete error: %v", err)
		return err
	}
	return nil
}

func (t *TokenPersistence) InvalidateRefreshToken(refreshToken string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"refresh_token": refreshToken}
	err := t.tokenDal.DeleteOne(ctx, filter)
	if err != nil {
		t.logger.Errorf("[InvalidateRefreshToken] delete error: %v", err)
		return err
	}
	return nil
}
