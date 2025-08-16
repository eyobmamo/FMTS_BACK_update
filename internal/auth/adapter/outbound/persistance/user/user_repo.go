package repo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"FMTS/internal/auth/domain/repository/user"
	dal "FMTS/internal/user/adapter/outbound/infra"
	model "FMTS/internal/user/domain/entity"
	"FMTS/pkg/utils"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UserAuthRepo struct {
	userDal dal.MongoDal[model.User, model.User]
	logger  utils.Logger
}

var _ repository.UserRepo = (*UserAuthRepo)(nil)

func NewUserAuthRepo(client *mongo.Client, dbName, collection string, logger utils.Logger) repository.UserRepo {
	userDal := dal.NewMongoDal[model.User, model.User](client, dbName, collection)
	return &UserAuthRepo{
		userDal: userDal,
		logger:  logger,
	}
}

func (u *UserAuthRepo) FindByEmail(email string) (*model.User, error) {
	filter := bson.M{"email": email, "is_deleted": false}
	projection := bson.M{}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user, err := u.userDal.FindOne(ctx, filter, projection)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		u.logger.Errorf("[FindByEmail] DB error: %v", err)
		return nil, err
	}
	return user, nil
}

func (u *UserAuthRepo) FindByID(id string) (*model.User, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %v", err)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}
	projection := bson.M{}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return u.userDal.FindOne(ctx, filter, projection)
}

func (u *UserAuthRepo) CreateUser(user model.User) (*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{"email": user.Email}
	projection := bson.M{"hashpassword": user.HashedPassword}
	userCreated, err := u.userDal.UpdateOne(ctx, filter, projection)
	if err != nil {
		u.logger.Errorf("[CreateUser] insert error: %v", err)
		return nil, err
	}
	return &userCreated, nil
}
