package util

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func SuccessResponse(w http.ResponseWriter, payload interface{}) {
	err := json.NewEncoder(w).Encode(payload)
	if err != nil {
		zap.L().Error("Writing response failed", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func NotFoundResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
}

func ErrorResponse(w http.ResponseWriter, err error, statusCode int) {
	w.WriteHeader(statusCode)
	err = json.NewEncoder(w).Encode(struct {
		Message string
		Code    int
	}{
		Message: err.Error(),
		Code:    statusCode,
	})
	if err != nil {
		zap.L().Error("Writing error response failed", zap.Int("code", statusCode), zap.Error(err))
	}
}

func GetStringVar(key string, r *http.Request) (string, error) {
	vars := mux.Vars(r)
	s, ok := vars[key]
	if !ok {
		return "", fmt.Errorf("missing request parameter %s", key)
	}
	return s, nil
}
