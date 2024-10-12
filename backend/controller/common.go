package controller

import (
	"encoding/json"
	"net/http"
)

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func errorResponse(w http.ResponseWriter, error error, status int) {
	errorResponse := ErrorResponse{
		Error: error.Error(),
	}
	res, _ := json.Marshal(errorResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(res)
}

func okResponse(w http.ResponseWriter, result []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func noContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

func okCreated(w http.ResponseWriter) {
	w.WriteHeader(http.StatusCreated)
}
