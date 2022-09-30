package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type JsonResponse struct {
	Error []string `json:"errors"`
	Data  any      `json:"data"`
}

func CustomResponse(w http.ResponseWriter, data any, errors []error, status int) (err error) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	errorFormated := make([]string, 0)
	for _, errorRaw := range errors {
		errorFormated = append(errorFormated, fmt.Sprintf("%s", errorRaw))
	}
	err = json.NewEncoder(w).Encode(JsonResponse{
		Error: errorFormated,
		Data:  data,
	})
	return err
}

func SuccessResponse(w http.ResponseWriter, data any, errors []error) (err error) {
	return CustomResponse(w, data, errors, http.StatusOK)
}
