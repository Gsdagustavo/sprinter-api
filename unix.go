//go:build unix

package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/Gsdagustavo/sprinter-api/domain/entities"
	"github.com/Gsdagustavo/sprinter-api/domain/entities/derr"
	"github.com/kardianos/service"
	"gopkg.in/yaml.v3"
)

func start() error {
	configsPath := os.Getenv("CONFIGS_PATH")
	if configsPath == "" {
		return errors.New("CONFIGS_PATH environment variable not set")
	}

	cfg, err := readCFGFile(configsPath)
	if err != nil {
		return derr.JoinError("failed to read config file", err)
	}

	file, err := configureOutput(cfg.LogSettings.LogDir)
	if err != nil {
		return derr.JoinError("failed to configure log outputs", err)
	}
	defer file.Close()

	s, err := newService(*cfg)
	if err != nil {
		return derr.JoinError("failed to create service", err)
	}

	err = s.Run()
	if err != nil {
		return derr.JoinError("failed to run service", err)
	}

	return nil
}

func (p *program) Start(s service.Service) error {
	slog.Info("received call to program#start")

	// Start should not block. Do the actual work async.
	go p.run(false)

	return nil
}

func (p *program) Stop(s service.Service) error {
	slog.Info("received call to program#stop")

	// Stop should not block. Return with a few seconds
	return nil
}

func configureOutput(logFolder string) (*os.File, error) {
	if logFolder == "" {
		return nil, nil
	}

	now := time.Now()
	logName := fmt.Sprintf("%s/%s.log", logFolder, now.Format("20060102150405"))

	err := os.MkdirAll(logFolder, os.ModePerm)
	if err != nil {
		return nil, derr.JoinError("failed to create log folder", err)
	}

	file, err := os.OpenFile(logName, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.WriteFile(logName, []byte(""), os.ModePerm)
			if err != nil {
				return nil, derr.JoinError("failed to write to file", err)
			}
		}

		return nil, derr.JoinError("failed to open log file", err)
	}

	log.SetOutput(io.MultiWriter(os.Stdout, file))
	return file, nil
}

func readCFGFile(cfgPath string) (*entities.Settings, error) {
	file, err := os.Open(cfgPath)
	if err != nil {
		return nil, derr.JoinError("failed to open file", err)
	}
	defer file.Close()

	b, err := io.ReadAll(file)
	if err != nil {
		return nil, derr.JoinError("failed to read file", err)
	}

	var cfg entities.Settings

	err = yaml.Unmarshal(b, &cfg)
	if err != nil {
		return nil, derr.JoinError("failed to decode file", err)
	}

	return &cfg, nil
}

func newService(cfg entities.Settings) (service.Service, error) {
	slog.Info("creating service")

	// Load the received arguments
	var args []string

	// Clean the executable arguments
	if len(os.Args) > 1 {
		for _, arg := range os.Args {
			if strings.Contains(arg, "configs") {
				args = append(args, arg)
			}
		}
	}

	svcConfig := &service.Config{
		Name:        "sprinter",
		DisplayName: "Sprinter API",
		Arguments:   args,
	}

	p := &program{
		cfg: cfg,
	}

	s, err := service.New(p, svcConfig)
	if err != nil {
		return nil, err
	}

	return s, nil
}
