package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	dal "FMTS/internal/user/adapter/outbound/infra"
	model "FMTS/internal/user/domain/entity"
	"FMTS/internal/user/port/outbound"
	"FMTS/pkg/utils"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	// "go.mongodb.org/mongo-driver/v2/mongo/options"
)

type UserPersistence struct {
	userDal dal.MongoDal[model.User, model.User]
	logger  utils.Logger
}

var _ outbound.UserRepoOutboundPort = (*UserPersistence)(nil)

func InitUserRepo(client *mongo.Client, dbName string, collection string, logger utils.Logger) outbound.UserRepoOutboundPort {
	userDal := dal.NewMongoDal[model.User, model.User](client, dbName, collection)
	return &UserPersistence{
		userDal: userDal,
		logger:  logger,
	}
}

// FindByEmailOrPhone checks if a user exists
func (u *UserPersistence) FindByEmailOrPhone(email string, phone string) (*model.User, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"email": email},
			{"phone": phone},
		},
		"is_deleted": false,
	}
	projection := bson.M{}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user, err := u.userDal.FindOne(ctx, filter, projection)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		u.logger.Errorf("[FindByEmailOrPhone] DB error: %v", err)
		return nil, err
	}
	return user, nil
}

// CreateUser inserts a user
func (u *UserPersistence) CreateUser(user model.User) (*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userCreated, err := u.userDal.InsertOne(ctx, user)
	if err != nil {
		u.logger.Errorf("[CreateUser] insert error: %v", err)
		return nil, err
	}
	return &userCreated, nil
}

// FindByID returns a user by ObjectID
func (u *UserPersistence) FindByID(id string) (*model.User, error) {
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

// FindAllUser returns all non-deleted users
func (u *UserPersistence) FindAllUser() ([]*model.User, error) {
	filter := bson.M{"is_deleted": false}
	projection := bson.M{}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return u.userDal.FindAll(ctx, filter, projection)
}

// UpdateUser updates mutable fields
func (u *UserPersistence) UpdateUser(user model.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": user.ID}
	update := bson.M{

		"updated_at": bson.NewDateTimeFromTime(time.Now()),
	}

	if user.FullName != "" {
		update["full_name"] = user.FullName
	}
	if user.FaydaID != "" {
		update["fayda_id"] = user.FullName
	}

	if user.PhoneNumber != "" {
		update["phone_number"] = user.FullName
	}

	if user.Email != "" {
		update["email"] = user.FullName
	}
	if user.CustomerType != "" {
		update["customer_type"] = user.FullName
	}

	_, err := u.userDal.UpdateOne(ctx, filter, update)
	if err != nil {
		u.logger.Errorf("[UpdateUser] update error: %v", err)
		return err
	}
	return nil
}

// UpdateSoftDelete sets is_deleted to true
func (u *UserPersistence) UpdateSoftDelete(id string) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid user ID: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = u.userDal.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		u.logger.Errorf("[UpdateSoftDelete] error: %v", err)
		return err
	}
	return nil
}
