package inbound

import "net/http"

type VehiclePortInterface interface {
	RegisterVehicle(w http.ResponseWriter, r *http.Request)
	GetVehicleByID(w http.ResponseWriter, r *http.Request)
	ListVehicles(w http.ResponseWriter, r *http.Request)
	UpdateVehicle(w http.ResponseWriter, r *http.Request)
	DeleteVehicle(w http.ResponseWriter, r *http.Request)
}
