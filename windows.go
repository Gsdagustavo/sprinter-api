//go:build windows

package main

import (
	_ "embed"
	"errors"
	"flag"
	"log/slog"
	"os"

	"github.com/Gsdagustavo/sprinter-api/domain/entities"
	"github.com/Gsdagustavo/sprinter-api/domain/entities/derr"
	"github.com/Gsdagustavo/sprinter-api/infrastructure/router/logger"
	"github.com/Gsdagustavo/sprinter-api/service"
	"gopkg.in/yaml.v3"
)

func Start() error {
	// Read the settings file
	cggPath, action := loadFlags()

	// Load the settings path
	if cggPath == "" {
		return errors.New("settings path is empty")
	}

	configsFile, err := os.Open(cggPath)
	if err != nil {
		return derr.JoinError("failed to open configs file", err)
	}
	defer configsFile.Close()

	var settings entities.Settings
	err = yaml.NewDecoder(configsFile).Decode(&settings)
	if err != nil {
		return derr.JoinError("failed to decode configs file", err)
	}

	err = logger.SetupLogger(settings)
	if err != nil {
		return derr.JoinError("failed to setup logger", err)
	}

	s, err := service.NewService(settings, cggPath)
	if err != nil {
		return derr.JoinError("failed to create service", err)
	}

	switch action {
	case "run", "":
		slog.Info("Run service")

		err = s.Run()
		if err != nil {
			return derr.JoinError("failed to run service", err)
		}

		return nil

	case "uninstall":
		slog.Info("Uninstall service")

		err = s.Uninstall()
		if err != nil {
			return derr.JoinError("failed to uninstall service", err)
		}

	case "install":
		slog.Info("Install service")

		err = s.Install()
		if err != nil {
			return derr.JoinError("failed to install service", err)
		}

	case "stop":
		slog.Info("Stop service")

		err = s.Stop()
		if err != nil {
			return derr.JoinError("failed to stop service", err)
		}
	}

	return nil
}

func loadFlags() (string, string) {
	var cfgPath string
	var action string

	// Try to read the configuration file from the command line arguments
	flag.StringVar(&cfgPath, "configs", "", "the path to the application config file")
	flag.StringVar(&action, "action", "", "the action to execute")
	flag.Parse()

	// If not provided as an argument, read from the environment
	if cfgPath == "" {
		cfgPath = os.Getenv("configs")
	}

	if cfgPath == "" {
		return "", ""
	}

	return cfgPath, action
}
