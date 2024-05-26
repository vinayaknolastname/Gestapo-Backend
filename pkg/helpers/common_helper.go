package helpers

import (
	"encoding/json"
	"net/http"
)

type jsonResponse struct {
	Status    bool       `json:"status"`
	Code      int        `json:"code,omitempty"`
	Message   string     `json:"message,omitempty"`
	Data      any        `json:"data,omitempty"`
	ErrorInfo *errorInfo `json:"error_info,omitempty"`
}
type errorInfo struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

func WriteJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(&jsonResponse{
		Status: true,
		Code:   status,
		Data:   data,
	})
}

func ErrorJson(w http.ResponseWriter, status int, err string) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	/*
		return json.NewEncoder(w).Encode(jsonResponse{
			Status: false,
			ErrorInfo: &errorInfo{
				StatusCode: status,
				Message:    err,
			},
		})
	*/
	return json.NewEncoder(w).Encode(jsonResponse{
		Code:    status,
		Status:  false,
		Message: err,
	})
}
