//go:build windows

package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/Gsdagustavo/sprinter-api/domain/entities"
	"github.com/kardianos/service"
)

func start() error {
	configsPath, action, terminal := loadFlags()
	cfg := readCFGFile(configsPath)

	if !terminal {
		file, err := configureOutput(cfg.LogDir)
		if err != nil {
			return errors.Join(errors.New("failed to configure log outputs"), err)
		}
		defer file.Close()
	}

	s, err := newService(*cfg)
	if err != nil {
		return errors.Join(errors.New("failed to create service"), err)
	}

	if strings.TrimSpace(action) == "" {
		return errors.New("action not provided")
	}

	switch action {
	case "run":
		slog.Info("run service")

		err = s.Run()
		if err != nil {
			return errors.Join(errors.New("failed to run service"), err)
		}
	case "uninstall":
		slog.Info("uninstalling service")

		err = s.Uninstall()
		if err != nil {
			return errors.Join(errors.New("failed to uninstall service"), err)
		}
	case "install":
		slog.Info("install service")

		err = s.Install()
		if err != nil {
			return errors.Join(errors.New("failed to install service"), err)
		}
	case "stop":
		log.Println("stopping service")

		err = s.Stop()
		if err != nil {
			return errors.Join(errors.New("failed to stop service"), err)
		}
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

func loadFlags() (string, string, bool) {
	var cfgPath string
	var action string
	var terminal bool

	// Try to read the configuration file from the command line arguments
	flag.StringVar(&cfgPath, "configs", "settings/dev-settings.toml", "the path to the application config file")
	flag.StringVar(&action, "action", "run", "the action to execute")
	flag.BoolVar(&terminal, "terminal", false, "display the logs in the terminal instead of in the log files")
	flag.Parse()

	// If not provided as argument, read from the environment
	if cfgPath == "" {
		cfgPath = os.Getenv("configs")
	}

	// And panic if the argument was not found
	if cfgPath == "" {
		panic("[configs] argument not found")
	}

	return cfgPath, action, terminal
}

func configureOutput(logFolder string) (*os.File, error) {
	if logFolder == "" {
		return nil, nil
	}

	now := time.Now()
	logName := fmt.Sprintf("%s/%s.log", logFolder, now.Format("20060102150405"))

	err := os.MkdirAll(logFolder, os.ModePerm)
	if err != nil {
		return nil, errors.Join(errors.New("failed to create log folder"))
	}

	file, err := os.OpenFile(logName, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.WriteFile(logName, []byte(""), os.ModePerm)
			if err != nil {
				return nil, errors.Join(errors.New("failed to create log file"))
			}
		}

		return nil, errors.Join(errors.New("failed to open log file"), err)
	}

	log.SetOutput(io.MultiWriter(os.Stdout, file))
	return file, nil
}

func readCFGFile(cfgPath string) *entities.Config {
	file, err := os.Open(cfgPath)
	if err != nil {
		panic(err)
	}

	b, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	err = file.Close()
	if err != nil {
		panic(err)
	}

	var cfg entities.Config

	_, err = toml.Decode(string(b), &cfg)
	if err != nil {
		panic(err)
	}

	return &cfg
}

func newService(cfg entities.Config) (service.Service, error) {
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
