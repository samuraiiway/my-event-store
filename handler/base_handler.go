package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"
)

func ParseBodyToMap(w http.ResponseWriter, r *http.Request) (map[string]interface{}, error) {
	body, readErr := ioutil.ReadAll(r.Body)
	if readErr != nil {
		http.Error(w, "bad_request", http.StatusBadRequest)
		return nil, readErr
	}

	result := map[string]interface{}{}
	jsonErr := json.Unmarshal(body, &result)
	if jsonErr != nil {
		http.Error(w, "bad_request", http.StatusBadRequest)
		return nil, jsonErr
	}

	return result, nil
}

func NewID() string {
	uid, _ := uuid.NewUUID()
	return uid.String()
}
