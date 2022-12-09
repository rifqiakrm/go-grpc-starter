package logger

import (
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// LogError define detailed error for logger
type LogError struct {
	Code    string
	Message string
	Error   error
	Jenjang string
	Detail  LogErrorDetail
}

// LogErrorDetail define detail of LogErrorDetail
type LogErrorDetail struct {
	Error string
}

// MarshalZerologObject function needed for logger
func (l LogError) MarshalZerologObject(e *zerolog.Event) {
	e.Str("code", l.Code).
		Str("message", l.Message).
		Str("jenjang", l.Jenjang).
		Str("error", l.Detail.Error)
}

type labelWarning struct {
	Warning string
}

type labelInfo struct {
	Info string
}

func (l labelWarning) MarshalZerologObject(e *zerolog.Event) {
	e.Str("warning", l.Warning)
}

func (l labelInfo) MarshalZerologObject(e *zerolog.Event) {
	e.Str("info", l.Info)
}

// Error logs an error message with the given error.
func Error(message string, detail LogError) {
	logEntry := log.Logger.With().
		Str("severity", "ERROR").
		Object("logging.googleapis.com/labels", labelInfo{Info: "applicationError"}).
		Str("logging.googleapis.com/trace", uuid.New().String()).
		Object("detail", detail).
		Str("message", message).
		Logger()
	logEntry.Info().Msg("")
}

// Warn logs a warning message with the given warning.
func Warn(err error) {
	logEntry := log.Logger.With().
		Str("severity", "WARNING").
		Object("logging.googleapis.com/labels", labelWarning{Warning: "applicationWarning"}).
		Str("logging.googleapis.com/trace", uuid.New().String()).
		Str("message", err.Error()).Logger()
	logEntry.Warn().Msg("")
}

// Info logs an info message with the given info.
func Info(message string) {
	logEntry := log.Logger.With().
		Str("severity", "INFO").
		Object("logging.googleapis.com/labels", labelInfo{Info: "applicationInfo"}).
		Str("logging.googleapis.com/trace", uuid.New().String()).
		Str("message", message).Logger()
	logEntry.Info().Msg("")
}
