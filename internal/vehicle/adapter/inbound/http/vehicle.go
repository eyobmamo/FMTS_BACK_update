package vehicle_handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	port "FMTS/internal/vehicle/port/inbound"
	context "FMTS/pkg/context"
	// "FMTS/pkg/utils"
	utility "FMTS/utils"

	// domain "FMTS/internal/user/domain/service"
	dto "FMTS/internal/vehicle/application"
	// model "FMTS/internal/vehicle/domain/entity"
	"FMTS/pkg/utils"
)

type VehicleHandler struct {
	vehicleService dto.VehicleService
	logger         utils.Logger
}

func NewVehicleHandler(service dto.VehicleService, logger utils.Logger) port.VehiclePortInterface {
	return &VehicleHandler{
		vehicleService: service,
		logger:         logger,
	}
}

func (h *VehicleHandler) RegisterVehicle(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateVehicleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Errorf("[RegisterVehicle] decode error: %v", err)
		utility.SendErrorResponse(w, "invalid request", http.StatusBadRequest, nil)
		return
	}

	if err := req.Validate(); err != nil {
		h.logger.Warnf("[RegisterVehicle] validation error: %v", err)
		utility.SendErrorResponse(w, err.Error(), http.StatusBadRequest, nil)
		return
	}

	ctx := context.ExtractUserContext(r)
	if ctx.UserID == "" {
		h.logger.Warnf("[RegisterVehicle] user context missing")
		utility.SendErrorResponse(w, "unauthorized", http.StatusUnauthorized, nil)
		return
	}

	created, err := h.vehicleService.RegisterVehicle(req, ctx.UserID)
	if err != nil {
		h.logger.Errorf("[RegisterVehicle] service error: %v", err)
		utility.SendErrorResponse(w, err.Error(), http.StatusInternalServerError, nil)
		return
	}

	utility.WriteSuccessResponse(w, created, "Vehicle registered successfully")
}

func (h *VehicleHandler) GetVehicleByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	result, err := h.vehicleService.GetVehicleByID(id)
	if err != nil {
		h.logger.Errorf("[GetVehicleByID] error: %v", err)
		utility.SendErrorResponse(w, err.Error(), http.StatusNotFound, nil)
		return
	}
	utility.WriteSuccessResponse(w, result, "Vehicle fetched successfully")
}

func (h *VehicleHandler) ListVehicles(w http.ResponseWriter, r *http.Request) {
	// User_id := chi.URLParam(r, "user_id")
	User_id := r.URL.Query().Get("user_id")

	vehicles, err := h.vehicleService.ListVehicles(User_id)
	if err != nil {
		h.logger.Errorf("[ListVehicles] error: %v", err)
		utility.SendErrorResponse(w, "failed to list vehicles", http.StatusInternalServerError, nil)
		return
	}
	utility.WriteSuccessResponse(w, vehicles, "Vehicles retrieved")
}

func (h *VehicleHandler) UpdateVehicle(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req dto.UpdateVehicleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Errorf("[UpdateVehicle] decode error: %v", err)
		utility.SendErrorResponse(w, "invalid input", http.StatusBadRequest, nil)
		return
	}

	if err := req.Validate(); err != nil {
		h.logger.Warnf("[UpdateVehicle] validation error: %v", err)
		utility.SendErrorResponse(w, err.Error(), http.StatusBadRequest, nil)
		return
	}

	updated, err := h.vehicleService.UpdateVehicle(id, req)
	if err != nil {
		h.logger.Errorf("[UpdateVehicle] update error: %v", err)
		utility.SendErrorResponse(w, err.Error(), http.StatusInternalServerError, nil)
		return
	}

	utility.WriteSuccessResponse(w, updated, "Vehicle updated successfully")
}

func (h *VehicleHandler) DeleteVehicle(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := h.vehicleService.DeleteVehicle(id)
	if err != nil {
		h.logger.Errorf("[DeleteVehicle] error: %v", err)
		utility.SendErrorResponse(w, err.Error(), http.StatusInternalServerError, nil)
		return
	}
	utility.WriteSuccessResponse(w, fmt.Sprintf("Vehicle ID: %s deleted successfully", id), "Vehicle deleted successfully")
}
