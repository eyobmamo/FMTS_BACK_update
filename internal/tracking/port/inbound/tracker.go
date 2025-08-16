package inbound

import "net/http"

type TrackerPortHandler interface {
	UpdateLocation(w http.ResponseWriter, r *http.Request)
	GetLetestViecleByViecleID(w http.ResponseWriter, r *http.Request)
	GetLetestLocationsOfViecleByUserID(w http.ResponseWriter, r *http.Request)
}
