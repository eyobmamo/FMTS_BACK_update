package inbound

import "net/http"

type TrackerPortHandler interface {
	UpdateLocation(w http.ResponseWriter, r *http.Request)
}
