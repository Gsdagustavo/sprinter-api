package router

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/VitorFranciscoDev/sprinter-api/domain/entities/derr"
	"github.com/VitorFranciscoDev/sprinter-api/domain/logger"
)

func Write(w http.ResponseWriter, v any) error {
	bytes, err := json.Marshal(v)
	if err != nil {
		return errors.New("failed to marshal response body")
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(bytes)
	if err != nil {
		return errors.New("failed to write response body")
	}

	return nil
}

func HandleError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	var clientErr derr.ClientError
	if errors.As(err, &clientErr) {
		w.WriteHeader(http.StatusBadRequest)
		response, e := json.Marshal(clientErr)
		if e != nil {
			slog.Error("failed to marshal client error response body", logger.Err(e))
			return
		}
		_, e = w.Write(response)
		if e != nil {
			slog.Error("failed to write client error response body", logger.Err(e))
		}
		return
	}

	var repositoryErr derr.RepositoryError
	if errors.As(err, &repositoryErr) {
		sanitized := derr.InternalServerError
		w.WriteHeader(http.StatusInternalServerError)
		response, e := json.Marshal(sanitized)
		if e != nil {
			slog.Error("failed to marshal repository error response body", logger.Err(e))
			return
		}
		_, e = w.Write(response)
		if e != nil {
			slog.Error("failed to write repository error response body", logger.Err(e))
		}
		return
	}

	sanitized := derr.InternalServerError
	w.WriteHeader(http.StatusInternalServerError)
	response, _ := json.Marshal(sanitized)
	_, e := w.Write(response)
	if e != nil {
		slog.Error("failed to write fallback error response body", logger.Err(e))
	}
}
