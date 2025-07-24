package vehicle

import (
	"context"
	"errors"
	// "fmt"

	// "fmt"
	"time"

	dal "FMTS/internal/user/adapter/outbound/infra"
	model "FMTS/internal/vehicle/domain/entity"
	vehicleOutboundPort "FMTS/internal/vehicle/port/outbound"
	"FMTS/pkg/utils"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type VehiclePersistence struct {
	vehicleDal dal.MongoDal[model.Vehicle, model.Vehicle]
	logger     utils.Logger
}

var _ vehicleOutboundPort.VehicleRepo = (*VehiclePersistence)(nil)

func InitVehicleRepo(client *mongo.Client, dbName string, collection string, logger utils.Logger) vehicleOutboundPort.VehicleRepo {
	vehicleDal := dal.NewMongoDal[model.Vehicle, model.Vehicle](client, dbName, collection)
	return &VehiclePersistence{
		vehicleDal: vehicleDal,
		logger:     logger,
	}
}

func (v *VehiclePersistence) FindByPlateNumber(plate string) (*model.Vehicle, error) {
	filter := bson.M{"plate_number": plate, "is_deleted": false}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	vehicle, err := v.vehicleDal.FindOne(ctx, filter, nil)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		v.logger.Errorf("[FindByPlateNumber] DB error: %v", err)
		return nil, err
	}
	return vehicle, nil
}

func (v *VehiclePersistence) CreateVehicle(vehicle model.Vehicle) (*model.Vehicle, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	createdVehicle, err := v.vehicleDal.InsertOne(ctx, vehicle)
	if err != nil {
		v.logger.Errorf("[CreateVehicle] insert error: %v", err)
		return nil, err
	}
	return &createdVehicle, nil
}

func (v *VehiclePersistence) FindByID(id string) (*model.Vehicle, error) {
	filter := bson.M{"_id": id, "is_deleted": false}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return v.vehicleDal.FindOne(ctx, filter, nil)
}

func (v *VehiclePersistence) FindAllVehicles(User_ID string) ([]*model.Vehicle, error) {
	filter := bson.M{"is_deleted": false}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// filter := bson.M{}
	if User_ID != "" {
		// objID, err := bson.ObjectIDFromHex(User_ID)
		// if err != nil {
		// 	v.logger.Infof("[find all vhicles (persistance) ] can't convert string id to object ID:  %v", err)
		// 	return nil, fmt.Errorf("INVALED_ID_PROVEDED")
		// }
		filter["owner_id"] = User_ID
	}
	projection := bson.M{}
	return v.vehicleDal.FindAll(ctx, filter, projection)
}

// func (v *VehiclePersistence) FindVehiclesWithFilter(query map[string]interface{}, page, limit int) ([]*model.Vehicle, error) {
// 	filter := bson.M{"is_deleted": false}
// 	for key, val := range query {
// 		filter[key] = val
// 	}

// 	skip := (page - 1) * limit
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	return v.vehicleDal.FindManyWithPagination(ctx, filter, skip, limit)
// }

func (v *VehiclePersistence) FindVehiclesWithFilter(query map[string]interface{}, page, limit int) ([]*model.Vehicle, error) {
	filter := bson.M{"is_deleted": false}

	// Merge query filters
	for key, val := range query {
		filter[key] = val
	}

	skip := int64((page - 1) * limit)
	limit64 := int64(limit)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	projection := bson.M{} // Optional: you can project specific fields if needed

	return v.vehicleDal.FindAllWithPagination(ctx, filter, projection, skip, limit64)
}

func (v *VehiclePersistence) UpdateVehicle(vehicle model.Vehicle) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": vehicle.ID}
	update := bson.M{"$set": vehicle}

	_, err := v.vehicleDal.UpdateOne(ctx, filter, update)
	if err != nil {
		v.logger.Errorf("[UpdateVehicle] update error: %v", err)
		return err
	}
	return nil
}

func (v *VehiclePersistence) UpdateSoftDelete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"is_deleted": true, "updated_at": time.Now()}}

	_, err := v.vehicleDal.UpdateOne(ctx, filter, update)
	if err != nil {
		v.logger.Errorf("[UpdateSoftDelete] error: %v", err)
		return err
	}
	return nil
}

// func (v *VehiclePersistence) SearchVehicles(search string, page, limit int) ([]*model.Vehicle, error) {
// 	filter := bson.M{"is_deleted": false, "$or": []bson.M{
// 		{"plate_number": bson.M{"$regex": search, "$options": "i"}},
// 		{"driver_name": bson.M{"$regex": search, "$options": "i"}},
// 	}}
// 	skip := (page - 1) * limit
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	return v.vehicleDal.FindManyWithPagination(ctx, filter, skip, limit)
// }

func (v *VehiclePersistence) SearchVehicles(search string, page, limit int) ([]*model.Vehicle, error) {
	filter := bson.M{
		"is_deleted": false,
		"$or": []bson.M{
			{"plate_number": bson.M{"$regex": search, "$options": "i"}},
			{"driver_name": bson.M{"$regex": search, "$options": "i"}},
		},
	}

	skip := int64((page - 1) * limit)
	limit64 := int64(limit)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	projection := bson.M{} // You can optionally include projections

	return v.vehicleDal.FindAllWithPagination(ctx, filter, projection, skip, limit64)
}
