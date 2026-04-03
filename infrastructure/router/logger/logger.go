package logger

import (
	"log/slog"
	"os"

	"github.com/takumakei/slogzap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLogger() {
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoderConfig.EncodeCaller = zapcore.FullCallerEncoder

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		zap.DebugLevel,
	)

	zapLogger := zap.New(core, zap.AddCaller())
	handler := slogzap.New(zapLogger)
	slog.SetDefault(slog.New(handler))
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
