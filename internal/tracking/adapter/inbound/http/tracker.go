// internal/tracking/handler/tracker_handler.go
package tracker

import (
	model "FMTS/internal/tracking/domain/entity"
	"FMTS/kafka"
	"encoding/json"
	"net/http"
)

type TrackerHandler struct {
	kafkaProducer *kafka.KafkaProducer
}

func NewTrackerHandler(producer *kafka.KafkaProducer) *TrackerHandler {
	return &TrackerHandler{
		kafkaProducer: producer,
	}
}

func (h *TrackerHandler) VehicleLocationUpdateHandler(w http.ResponseWriter, r *http.Request) {
	var loc model.VehicleLocation
	if err := json.NewDecoder(r.Body).Decode(&loc); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	err := h.kafkaProducer.ProduceVehicleLocation(r.Context(), loc)
	if err != nil {
		http.Error(w, "failed to publish location", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
