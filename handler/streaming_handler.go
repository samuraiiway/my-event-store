package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/samuraiiway/my-event-store/service"
)

const (
	AGGREGATED_STREAM_LISTENER_PATH = "/streaming/aggregation/{namespace}/{aggregated_id}"
)

func AggregatedStreamListener(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespace := vars["namespace"]
	aggregatedID := vars["aggregated_id"]
	clientID, clientOk := vars["client_id"]
	if !clientOk {
		clientID = NewID()
	}

	flusher, flusherOk := w.(http.Flusher)

	if !flusherOk {
		fmt.Println("Unsupported steaming")
		http.Error(w, "Unsupported steaming", http.StatusInternalServerError)
		return
	}

	endSignal := r.Context().Done()
	ch := service.RegisterClientListener(namespace, aggregatedID, clientID)

	go func() {
		<-endSignal
		service.DeregisterClientListener(namespace, aggregatedID, clientID)
	}()

	w.Header().Set("Content-Type", "text/event-stream")

	for {
		data, ok := <-ch
		if !ok {
			return
		}
		response, _ := json.Marshal(data)
		fmt.Fprintf(w, "data: %s\n\n", response)
		flusher.Flush()
	}
}
