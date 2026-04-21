//go:build unix

package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Gsdagustavo/sprinter-api/domain/entities"
	"github.com/Gsdagustavo/sprinter-api/infrastructure"
	"github.com/Gsdagustavo/sprinter-api/infrastructure/router/logger"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func Start() error {
	settings, err := loadSettingsFromEnv()
	if err != nil {
		return errors.Join(errors.New("failed to load settings from environment variables"), err)
	}

	// Order the modules
	router := mux.NewRouter()

	err = infrastructure.SetupModules(router, settings)
	if err != nil {
		return errors.Join(errors.New("failed to setup infrastructure"), err)
	}

	err = logger.SetupLogger(settings)
	if err != nil {
		return errors.Join(errors.New("failed to setup logger"), err)
	}

	srv := &http.Server{
		Handler: handlers.CORS(
			handlers.AllowedOrigins(settings.CORSConfig.CORSOrigins),
			handlers.AllowedHeaders([]string{"Authorization", "Content-Type", "Accept"}),
			handlers.AllowedMethods([]string{"GET", "POST", "DELETE", "UPDATE"}),
			handlers.AllowCredentials(),
		)(router),
		Addr: fmt.Sprintf(":%d", settings.Server.Port),

		WriteTimeout: 1200 * time.Second,
		ReadTimeout:  1200 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())

	return nil
}

func requireEnv(key string) (string, error) {
	v := os.Getenv(key)
	if v == "" {
		return "", errors.Join(fmt.Errorf("env %q is required but not set", key))
	}
	return v, nil
}

func optionalEnv(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}

func requireEnvInt(key string) (int, error) {
	raw, err := requireEnv(key)
	if err != nil {
		return 0, errors.Join(errors.New("failed to read env as int"), err)
	}
	v, err := strconv.Atoi(raw)
	if err != nil {
		return 0, errors.Join(fmt.Errorf("env %q must be a valid integer", key), err)
	}
	return v, nil
}

func optionalEnvInt(key string, defaultValue int) (int, error) {
	raw := os.Getenv(key)
	if raw == "" {
		return defaultValue, nil
	}
	v, err := strconv.Atoi(raw)
	if err != nil {
		return 0, errors.Join(fmt.Errorf("env %q must be a valid integer", key), err)
	}
	return v, nil
}

func optionalEnvDuration(key string, defaultValue time.Duration) (time.Duration, error) {
	raw := os.Getenv(key)
	if raw == "" {
		return defaultValue, nil
	}
	d, err := time.ParseDuration(raw)
	if err != nil {
		return 0, errors.Join(fmt.Errorf("env %q must be a valid duration (e.g. 30s, 1m)", key), err)
	}
	return d, nil
}

func loadSettingsFromEnv() (s entities.Settings, err error) {
	// ── Server ───────────────────────────────────────────────────────────────

	s.Server.Host, err = requireEnv("SERVER_HOST")
	if err != nil {
		return s, err
	}

	s.Server.Port, err = requireEnvInt("SERVER_PORT")
	if err != nil {
		return s, err
	}

	s.Server.Domain, err = requireEnv("SERVER_DOMAIN")
	if err != nil {
		return s, err
	}

	// ── Database ─────────────────────────────────────────────────────────────

	s.RepositorySettings.Host, err = requireEnv("DB_HOST")
	if err != nil {
		return s, err
	}

	s.RepositorySettings.Port, err = requireEnv("DB_PORT")
	if err != nil {
		return s, err
	}

	s.RepositorySettings.Name, err = requireEnv("DB_NAME")
	if err != nil {
		return s, err
	}

	s.RepositorySettings.User, err = requireEnv("DB_USER")
	if err != nil {
		return s, err
	}

	s.RepositorySettings.Password, err = requireEnv("DB_PASSWORD")
	if err != nil {
		return s, err
	}

	// ── Service ───────────────────────────────────────────────────────────────

	s.ServiceConfig.Name, err = requireEnv("SERVICE_NAME")
	if err != nil {
		return s, err
	}

	s.ServiceConfig.Display, err = requireEnv("SERVICE_DISPLAY")
	if err != nil {
		return s, err
	}

	s.ServiceConfig.Description, err = requireEnv("SERVICE_DESCRIPTION")
	if err != nil {
		return s, err
	}

	// ── Environment ───────────────────────────────────────────────────────────

	envType, err := requireEnv("ENVIRONMENT_TYPE")
	if err != nil {
		return s, err
	}
	s.EnvironmentSettings.EnvironmentType = entities.EnvironmentType(envType)

	// ── Logging ───────────────────────────────────────────────────────────────

	s.LogSettings.LogDir = optionalEnv("LOG_DIR", "./logs")

	// ── PASETO ────────────────────────────────────────────────────────────────

	s.PasetoSettings.SecurityKey, err = requireEnv("PASETO_SECURITY_KEY")
	if err != nil {
		return s, err
	}

	// ── CORS ──────────────────────────────────────────────────────────────────

	rawOrigins := optionalEnv("CORS_ORIGINS", "*")
	s.CORSConfig.CORSOrigins = strings.Split(rawOrigins, ",")

	// ── File Storage ──────────────────────────────────────────────────────────

	s.FileStorageSettings.StorageFolder = optionalEnv("FILE_STORAGE_FOLDER", "./storage")

	// ── SMTP ──────────────────────────────────────────────────────────────────

	s.SMTPSettings.Host, err = requireEnv("SMTP_HOST")
	if err != nil {
		return s, err
	}

	s.SMTPSettings.Port, err = requireEnvInt("SMTP_PORT")
	if err != nil {
		return s, err
	}

	s.SMTPSettings.User, err = requireEnv("SMTP_USER")
	if err != nil {
		return s, err
	}

	s.SMTPSettings.Password, err = requireEnv("SMTP_PASSWORD")
	if err != nil {
		return s, err
	}

	s.SMTPSettings.From, err = requireEnv("SMTP_FROM")
	if err != nil {
		return s, err
	}

	s.SMTPSettings.MaxConnections, err = optionalEnvInt("SMTP_MAX_CONNECTIONS", 5)
	if err != nil {
		return s, err
	}

	s.SMTPSettings.IdleTimeout, err = optionalEnvDuration("SMTP_IDLE_TIMEOUT", 30*time.Second)
	if err != nil {
		return s, err
	}

	s.SMTPSettings.PoolWaitTimeout, err = optionalEnvDuration("SMTP_POOL_WAIT_TIMEOUT", 10*time.Second)
	if err != nil {
		return s, err
	}

	return s, nil
}
