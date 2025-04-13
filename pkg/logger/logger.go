package logger

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	requestIDKey = "request_id"
)

// Logger представляет логгер приложения
type Logger struct {
	*zap.Logger
}

// NewLogger создает новый экземпляр логгера
func NewLogger(level string) (*Logger, error) {
	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(level)); err != nil {
		return nil, fmt.Errorf("невозможно распарсить уровень логирования: %w", err)
	}

	config := zap.Config{
		Level:            zap.NewAtomicLevelAt(zapLevel),
		Development:      false,
		Encoding:         "json",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	zapLogger, err := config.Build()
	if err != nil {
		return nil, fmt.Errorf("ошибка инициализации логгера: %w", err)
	}

	return &Logger{zapLogger}, nil
}

// WithRequestID добавляет request_id в логгер из контекста
func (l *Logger) WithRequestID(ctx context.Context) *Logger {
	requestID, ok := ctx.Value(requestIDKey).(string)
	if !ok || requestID == "" {
		requestID = uuid.New().String()
	}

	return &Logger{l.With(zap.String(requestIDKey, requestID))}
}

// Error создает поле ошибки для логгера
func Error(err error) zap.Field {
	return zap.Error(err)
}

// String создает строковое поле для логгера
func String(key, value string) zap.Field {
	return zap.String(key, value)
}

// Int создает целочисленное поле для логгера
func Int(key string, value int) zap.Field {
	return zap.Int(key, value)
}

// Any создает поле произвольного типа для логгера
func Any(key string, value interface{}) zap.Field {
	return zap.Any(key, value)
}
