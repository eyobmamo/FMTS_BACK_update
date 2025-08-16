package persistence

import (
	entity "FMTS/internal/tracking/domain/entity"
	// port "FMTS/internal/tracking/port/outbound"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TimescaleTrackerRepo struct {
	db *pgxpool.Pool
}

func NewTimescaleTrackerRepo(db *pgxpool.Pool) *TimescaleTrackerRepo {
	return &TimescaleTrackerRepo{db: db}
}

func (r *TimescaleTrackerRepo) UpdateLocation(ctx context.Context, location entity.VehicleLocation) (entity.VehicleLocation, error) {
	query := `
		INSERT INTO vehicle_locations (owner_id, vehicle_id, latitude, longitude, speed, timestamp)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id;
	`

	_, err := r.db.Exec(ctx, query,
		location.OwnerID,
		location.VehicleID,
		location.Latitude,
		location.Longitude,
		location.Speed,
		location.Timestamp,
	)

	if err != nil {
		return entity.VehicleLocation{}, fmt.Errorf("failed to insert location: %w", err)
	}

	return location, nil
}

func (r *TimescaleTrackerRepo) GetLatestVehicleLocationByID(ctx context.Context, vehicleID string) (entity.VehicleLocation, error) {
	fmt.Printf("vehile ID : %v", vehicleID)
	const query = `
        SELECT owner_id, vehicle_id, latitude, longitude, speed, timestamp
        FROM vehicle_locations
        WHERE vehicle_id = $1
        ORDER BY timestamp DESC
        LIMIT 1;
    `

	var loc entity.VehicleLocation
	err := r.db.QueryRow(ctx, query, "6881554e4f64d5fd1c00ba5b").Scan(
		&loc.OwnerID,
		&loc.VehicleID,
		&loc.Latitude,
		&loc.Longitude,
		&loc.Speed,
		&loc.Timestamp,
	)
	if err != nil {
		return entity.VehicleLocation{}, fmt.Errorf("failed to get latest location: %w", err)
	}

	return loc, nil
}
func (r *TimescaleTrackerRepo) GetLatestVehicleLocationsByUserID(ctx context.Context, userID string) ([]*entity.VehicleLocation, error) {
	const query = `
		SELECT DISTINCT ON (vehicle_id) owner_id, vehicle_id, latitude, longitude, speed, timestamp
		FROM vehicle_locations
		WHERE owner_id = $1
		ORDER BY vehicle_id, timestamp DESC;
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query latest locations by user: %w", err)
	}
	defer rows.Close()

	var locations []*entity.VehicleLocation
	for rows.Next() {
		var loc entity.VehicleLocation
		if err := rows.Scan(
			&loc.OwnerID,
			&loc.VehicleID,
			&loc.Latitude,
			&loc.Longitude,
			&loc.Speed,
			&loc.Timestamp,
		); err != nil {
			return nil, fmt.Errorf("failed to scan location: %w", err)
		}
		locations = append(locations, &loc)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return locations, nil
}
