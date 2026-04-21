package logger

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/Gsdagustavo/sprinter-api/domain/entities"
	slogzap "github.com/samber/slog-zap/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ctxKey string

const (
	slogFields      ctxKey = "slog_fields"
	logNameTemplate        = "20060102"
)

// AddLogValueToContext adds a slog attribute to the provided context so that it will be
// included in any Record created with such context
func AddLogValueToContext(parent context.Context, attr ...slog.Attr) context.Context {
	if parent == nil {
		parent = context.Background()
	}

	if v, ok := parent.Value(slogFields).([]slog.Attr); ok {
		for _, i := range attr {
			v = append(v, i)
		}
		return context.WithValue(parent, slogFields, v)
	}

	var v []slog.Attr
	for _, i := range attr {
		v = append(v, i)
	}
	return context.WithValue(parent, slogFields, v)
}

// GetLogValuesFromContext retrieves the slice with all log values added to ctx using AddLogValueToContext
func GetLogValuesFromContext(ctx context.Context) []slog.Attr {
	attrs, ok := ctx.Value(slogFields).([]slog.Attr)
	if ok {
		return attrs
	}
	return nil
}

// CreateLogHandler creates a properly configured slog.Handler using structured logging and the Zap library
func CreateLogHandler(outPath string, settings entities.Settings) (slog.Handler, error) {
	// Customize the logger according to the environment where the server is running
	switch {
	case settings.IsProduction():
		zapLoggerConfig := zap.NewProductionConfig()
		if outPath != "" {
			zapLoggerConfig.OutputPaths = []string{outPath}
		}

		zapLogger, err := zapLoggerConfig.Build()
		if err != nil {
			return nil, errors.Join(errors.New("failed to build production logger"), err)
		}

		zapOptions := slogzap.Option{
			Level:     slog.LevelInfo,
			Logger:    zapLogger,
			AddSource: true,
			AttrFromContext: []func(ctx context.Context) []slog.Attr{
				GetLogValuesFromContext,
			},
		}

		return zapOptions.NewZapHandler(), nil
	case settings.IsLocal():
		consoleEncoderConfig := zap.NewDevelopmentEncoderConfig()
		consoleEncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		consoleEncoder := zapcore.NewConsoleEncoder(consoleEncoderConfig)

		fileEncoderConfig := zap.NewDevelopmentEncoderConfig()
		fileEncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		fileEncoder := zapcore.NewConsoleEncoder(fileEncoderConfig)

		cores := []zapcore.Core{
			zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel),
		}

		if outPath != "" {
			f, err := os.OpenFile(outPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return nil, fmt.Errorf("failed to open log file: %w", err)
			}

			cores = append(cores, zapcore.NewCore(fileEncoder, zapcore.AddSync(f), zapcore.DebugLevel))
		}

		zapLogger := zap.New(zapcore.NewTee(cores...), zap.AddCaller())

		zapOptions := slogzap.Option{
			Level:     slog.LevelDebug,
			Logger:    zapLogger,
			AddSource: true,
			AttrFromContext: []func(ctx context.Context) []slog.Attr{
				GetLogValuesFromContext,
			},
		}

		return zapOptions.NewZapHandler(), nil

	default:
		return nil, errors.New("unknown environment: " + string(settings.EnvironmentSettings.EnvironmentType))
	}
}

func SetupLogger(settings entities.Settings) error {
	serverLogs, httpLogs := selectLogOutput(settings)

	if serverLogs != "" {
		err := os.MkdirAll(filepath.Dir(serverLogs), 0755)
		if err != nil {
			return errors.Join(errors.New("failed to create log dir"), err)
		}
	}

	if httpLogs != "" {
		err := os.MkdirAll(filepath.Dir(httpLogs), 0755)
		if err != nil {
			return errors.Join(errors.New("failed to create http log dir"), err)
		}
	}

	// Initialize the log handler
	logHandler, err := CreateLogHandler(serverLogs, settings)
	if err != nil {
		return errors.Join(errors.New("failed to create log handler"), err)
	}

	// Wrap with the context handler to always process the context variables in the log.
	slog.SetDefault(slog.New(logHandler))

	return nil
}

// Returns the path to the application and middleware log files, and whether the log file rotation should be enabled
func selectLogOutput(settings entities.Settings) (string, string) {
	logDir := settings.LogSettings.LogDir
	if logDir == "" {
		return "", ""
	}

	formattedDate := time.Now().Format(logNameTemplate)
	applicationLog := filepath.Join(logDir, fmt.Sprintf("%s.log", formattedDate))
	middlewareLog := filepath.Join(logDir, fmt.Sprintf("%s-http.log", formattedDate))

	return applicationLog, middlewareLog
}

// Err returns an Attr for an error value
func Err(err error) slog.Attr {
	if err == nil {
		return slog.Attr{
			Key:   "cause",
			Value: slog.AnyValue(err),
		}
	}

	return slog.Attr{
		Key:   "cause",
		Value: slog.StringValue(err.Error()),
	}
}
