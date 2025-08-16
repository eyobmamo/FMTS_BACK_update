// internal/tracking/handler/tracker_handler.go
package tracker

import (
	app "FMTS/internal/tracking/application"
	model "FMTS/internal/tracking/domain/entity"
	"FMTS/kafka"
	"FMTS/pkg/utils"
	utility "FMTS/utils"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type TrackerHandler struct {
	kafkaProducer *kafka.KafkaProducer
	AppTracker    app.TrackerApplication
	logger        utils.Logger
}

func NewTrackerHandler(producer *kafka.KafkaProducer, appTracker app.TrackerApplication, logger utility.Logger) *TrackerHandler {
	return &TrackerHandler{
		kafkaProducer: producer,
		AppTracker:    appTracker,
		logger:        logger,
	}
}

func (h *TrackerHandler) UpdateLocation(w http.ResponseWriter, r *http.Request) {
	var req model.VehicleLocation
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Errorf("[UpdateLocation] failed to decode request: %v", err)
		utility.SendErrorResponse(w, "invalid request format", http.StatusBadRequest, nil)
		return
	}
	if err := req.Validate(); err != nil {
		h.logger.Warnf("[UpdateLocation] validation failed: %v", err)
		utility.SendErrorResponse(w, err.Error(), http.StatusBadRequest, nil)
		return
	}

	locationUpdated, err := h.AppTracker.UpdateLocation(r.Context(), req)
	if err != nil {
		h.logger.Errorf("[UpdateLocation] service error: %v", err)
		utility.SendErrorResponse(w, err.Error(), http.StatusInternalServerError, nil)
		return
	}
	utility.WriteSuccessResponse(w, locationUpdated, "Location updated successfully")
}

func (h *TrackerHandler) GetLetestViecleByViecleID(w http.ResponseWriter, r *http.Request) {
	vehicleID := chi.URLParam(r, "vehicle_id")
	location, err := h.AppTracker.GetLatestVehicleLocationByID(r.Context(), vehicleID)
	if err != nil {
		h.logger.Errorf("[GetLatestVehicleByVehicleID] failed: %v", err)
		utility.SendErrorResponse(w, err.Error(), http.StatusNotFound, nil)
		return
	}
	utility.WriteSuccessResponse(w, location, "Latest vehicle location fetched successfully")
}

func (h *TrackerHandler) GetLetestLocationsOfViecleByUserID(w http.ResponseWriter, r *http.Request) {
	// userID, ok := r.Context().Value("user_id").(string)
	// if !ok || userID == "" {
	// 	h.logger.Warnf("[GetLetestLocationsOfViecleByUserID] user_id not found in context")
	// 	utility.SendErrorResponse(w, "user_id not found in context", http.StatusBadRequest, nil)
	// 	return
	// }

	userID := "687ab88a9a3d33a856cd5000"
	fmt.Printf("user id: %v", userID)
	locations, err := h.AppTracker.GetLatestVehicleLocationsByUserID(r.Context(), userID)
	if err != nil {
		h.logger.Errorf("[GetLetestLocationsOfViecleByUserID] failed: %v", err)
		utility.SendErrorResponse(w, err.Error(), http.StatusNotFound, nil)
		return
	}
	utility.WriteSuccessResponse(w, locations, "Latest vehicle locations fetched successfully")
}
