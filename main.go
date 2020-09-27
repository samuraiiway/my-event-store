package main

import (
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/gorilla/mux"
	"github.com/samuraiiway/my-event-store/handler"
	"github.com/samuraiiway/my-event-store/service"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc(handler.CREATE_EVENT_PATH, handler.CreateEventHandler).Methods("POST")
	router.HandleFunc(handler.CREATE_AGGREGATED_PATH, handler.CreateAggregatedHandler).Methods("POST")
	router.HandleFunc(handler.AGGREGATED_STREAM_LISTENER_PATH, handler.AggregatedStreamListener).Methods("GET")

	go monitoring()
	http.ListenAndServe(":8000", router)
}

func monitoring() {
	for {
		time.Sleep(4 * time.Second)
		fmt.Printf("========== Monitoring : %v ==========\n", time.Now())
		fmt.Printf("Thread : %v\n", runtime.NumGoroutine())
		fmt.Printf("Root Aggregated Channels : %v\n", service.GetRootAggregatedChannel())
		fmt.Printf("Root Aggregated Listeners : %v\n", service.GetRootAggregatedListener())
		fmt.Printf("Root Tasks : %v\n\n", service.GetRootTasks())
	}
}
