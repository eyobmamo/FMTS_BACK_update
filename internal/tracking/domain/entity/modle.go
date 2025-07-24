package models

import (
	"time"
)

type VehicleLocation struct {
	OwnerID   string  `json:"owner_id" bson:"owner_id"`               // ID of the owner/user who owns the vehicle
	VehicleID string  `json:"vehicle_id" bson:"vehicle_id"`           // Unique vehicle identifier
	Latitude  float64 `json:"latitude" bson:"latitude"`               // Current latitude
	Longitude float64 `json:"longitude" bson:"longitude"`             // Current longitude
	Speed     float64 `json:"speed,omitempty" bson:"speed,omitempty"` // Optional: speed in km/h
	// Heading   float64   `json:"heading,omitempty" bson:"heading,omitempty"`     // Optional: compass heading
	Timestamp time.Time `json:"timestamp" bson:"timestamp"` // UTC time of the location record
}
