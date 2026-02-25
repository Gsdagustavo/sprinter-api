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

func WriteInternalError(w http.ResponseWriter) {
	http.Error(w, "Internal server error", http.StatusInternalServerError)
}

func WriteBadRequest(w http.ResponseWriter) {
	http.Error(w, "Bad request", http.StatusBadRequest)
}

func WriteUnauthorized(w http.ResponseWriter) {
	http.Error(w, "Unauthorized", http.StatusUnauthorized)
}

func WriteForbidden(w http.ResponseWriter) {
	http.Error(w, "Forbidden", http.StatusForbidden)
}

func HandleError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	response, err := json.Marshal(err)
	if err != nil {
		slog.Error("failed to marshal error response body", logger.Err(err))
		return
	}

	var clientErr derr.ClientError
	if errors.As(err, &clientErr) {
		w.WriteHeader(http.StatusBadRequest)
		_, e := w.Write(response)
		if e != nil {
			slog.Error("failed to write client error response body", logger.Err(e))
		}
		return
	}

	var repositoryError derr.RepositoryError
	if errors.As(err, &repositoryError) {
		w.WriteHeader(http.StatusInternalServerError)
		_, e := w.Write(response)
		if e != nil {
			slog.Error("failed to write repository error response body", logger.Err(e))
		}
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
	_, e := w.Write(response)
	if e != nil {
		slog.Error("failed to write repository error response body", logger.Err(e))
	}
}
