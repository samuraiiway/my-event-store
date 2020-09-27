package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/samuraiiway/my-event-store/service"
)

const (
	CREATE_EVENT_PATH = "/event/{namespace}"
)

func CreateEventHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespace := vars["namespace"]
	body, err := ParseBodyToMap(w, r)
	if err != nil {
		return
	}

	body["id"] = NewID()
	service.SendTask(namespace, body)

	response, _ := json.Marshal(body)

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
