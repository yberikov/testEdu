package handler

import (
	"fmt"
	"net/http"
	"strconv"
)

func GetDevice(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Invalid Method: should be GET", http.StatusBadRequest)
		return
	}
	id := r.URL.Query().Get("id")

	if id == "" {
		http.Error(w, "Bad Request: Missing 'id' parameter", http.StatusBadRequest)
		return
	}

	_, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Bad Request: 'id' must be a number", http.StatusBadRequest)
		return
	}

	response := fmt.Sprintf("get device with id %s", id)
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(response))
	if err != nil {
		return
	}
}
