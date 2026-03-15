package logger

import (
	"log/slog"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLogger() error {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:          "time",
		EncodeTime:       zapcore.ISO8601TimeEncoder,
		EncodeLevel:      zapcore.CapitalColorLevelEncoder,
		ConsoleSeparator: "  ",
	}

	encoder := zapcore.NewConsoleEncoder(encoderConfig)
	core := zapcore.NewCore(encoder, os.Stdout, zapcore.DebugLevel)
	logger := zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(logger)
	return nil
}

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
