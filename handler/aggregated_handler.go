package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/samuraiiway/my-event-store/model"
	"github.com/samuraiiway/my-event-store/service"
)

const (
	CREATE_AGGREGATED_PATH = "/aggregation/{namespace}"
)

func CreateAggregatedHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespace := vars["namespace"]
	var body model.AggregatedCreateRequest

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "bad_request", http.StatusBadRequest)
		return
	}

	if body.AggregatedID == "" {
		body.AggregatedID = NewID()
	}

	service.NewAggregatedTask(namespace, &body)
	service.NewAggregatedChannel(namespace)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&body)
}
