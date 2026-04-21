package service

import (
	_ "embed"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/Gsdagustavo/sprinter-api/domain/entities"
	"github.com/Gsdagustavo/sprinter-api/infrastructure"
	"github.com/Gsdagustavo/sprinter-api/infrastructure/router/logger"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/kardianos/service"
)

func NewService(settings entities.Settings, settingsPath string) (service.Service, error) {
	slog.Info("Creating service")

	serviceConfig := settings.ServiceConfig
	if serviceConfig.Name == "" {
		return nil, errors.New("service name is empty")
	}

	// Service display
	if serviceConfig.Display == "" {
		return nil, errors.New("service display name is empty")
	}

	// Load the received arguments
	args := []string{"--configs", settingsPath}
	svcConfig := &service.Config{
		Name:        serviceConfig.Name,
		DisplayName: serviceConfig.Display,
		Description: serviceConfig.Description,
		Arguments:   args,
	}

	prg := &program{settings: settings}

	s, err := service.New(prg, svcConfig)
	if err != nil {
		return nil, errors.Join(errors.New("failed to create service"), err)
	}

	_, err = s.Logger(nil)
	if err != nil {
		return nil, errors.Join(errors.New("failed to set logger"), err)
	}

	return s, nil
}

type program struct {
	settings entities.Settings
}

func (p *program) Start(_ service.Service) error {
	// Start should not block. Do the actual work async.
	go p.run()

	return nil
}

func (p *program) Stop(_ service.Service) error {
	slog.Info("Stopping server")

	// Stop should not block. Return with a few seconds.
	return nil
}

func (p *program) run() {
	settings := p.settings

	router := mux.NewRouter()
	err := infrastructure.SetupModules(router, settings)
	if err != nil {
		slog.Error("failed to setup modules", logger.Err(err))
		return
	}

	handler := handlers.CORS(
		handlers.AllowedOrigins(settings.CORSConfig.CORSOrigins),
		handlers.AllowedHeaders([]string{"Authorization", "Content-Type", "Accept"}),
		handlers.AllowedMethods([]string{"GET", "POST", "DELETE", "UPDATE"}),
		handlers.AllowCredentials(),
	)(router)

	srv := &http.Server{
		Handler:      handler,
		Addr:         fmt.Sprintf(":%d", settings.Server.Port),
		WriteTimeout: 120 * time.Second,
		ReadTimeout:  120 * time.Second,
	}

	slog.Info("Starting server", slog.Int("port", settings.Server.Port))
	log.Fatal(srv.ListenAndServe())
}
