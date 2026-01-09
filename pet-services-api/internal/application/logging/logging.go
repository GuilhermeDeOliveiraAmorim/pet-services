package logging

import (
	"context"
	"log/slog"
	"os"

	"pet-services-api/internal/application/exceptions"
	"github.com/lmittmann/tint"
)

// Logger retorna o logger interno como *slog.Logger
func (s *SlogLogger) Logger() *slog.Logger {
	return s.logger
}

const TimeFormat = "2006-01-02 15:04:05"

type Layer struct {
	ENTITY                                     string
	FACTORIES                                  string
	INFRASTRUCTURE_REPOSITORIES_IMPLEMENTATION string
	INTERFACE_HANDLERS                         string
	USECASES                                   string
	CONFIGURATION                              string
	MIDDLEWARES                                string
	SERVICES                                   string
	SERVER                                     string
}

type TypeLog struct {
	ERROR   string
	INFO    string
	WARNING string
}

type DefaultMessages struct {
	START string
	END   string
}

var LoggerLayers = Layer{
	ENTITY:    "ENTITY",
	FACTORIES: "FACTORIES",
	INFRASTRUCTURE_REPOSITORIES_IMPLEMENTATION: "INFRASTRUCTURE_REPOSITORIES_IMPLEMENTATION",
	INTERFACE_HANDLERS:                         "INTERFACE_HANDLERS",
	USECASES:                                   "USECASES",
	CONFIGURATION:                              "CONFIGURATION",
	MIDDLEWARES:                                "MIDDLEWARES",
	SERVICES:                                   "SERVICES",
	SERVER:                                     "SERVER",
}

var LoggerTypes = TypeLog{
	ERROR:   "ERROR",
	INFO:    "INFO",
	WARNING: "WARNING",
}

var DEFAULTMESSAGES = DefaultMessages{
	START: "Iniciando operação",
	END:   "Operação concluída",
}

type Logger struct {
	Context  context.Context             `json:"context"`
	Code     int                         `json:"code"`
	Message  string                      `json:"message"`
	From     string                      `json:"from"`
	Layer    string                      `json:"layer"`
	TypeLog  string                      `json:"type_log"`
	Error    error                       `json:"error"`
	Problems []exceptions.ProblemDetails `json:"problems"`
}

type LoggerService interface {
	Log(log Logger)
}

type SlogLogger struct {
	logger *slog.Logger
}

func NewSlogLogger() *SlogLogger {
	return &SlogLogger{
		logger: slog.New(tint.NewHandler(os.Stdout, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: TimeFormat,
		})),
	}
}

func (s *SlogLogger) Log(log Logger) {
	switch log.TypeLog {
	case "ERROR":
		s.logger.ErrorContext(
			log.Context,
			"ERROR",
			"code:", log.Code,
			"message:", log.Message,
			"from:", log.From,
			"layer:", log.Layer,
			"error:", log.Error,
			"problems:", log.Problems,
		)
	case "INFO":
		s.logger.InfoContext(
			log.Context,
			"INFO",
			"code:", log.Code,
			"message:", log.Message,
			"from:", log.From,
			"layer:", log.Layer,
			"error:", log.Error,
			"problems:", log.Problems,
		)
	case "WARNING":
		s.logger.WarnContext(
			log.Context,
			"WARNING",
			"code:", log.Code,
			"message:", log.Message,
			"from:", log.From,
			"layer:", log.Layer,
			"error:", log.Error,
			"problems:", log.Problems,
		)
	}
}

func LogWithProblem(logger LoggerService, ctx context.Context, from, layer, title string, err error, code int, errorType exceptions.ErrorType) []exceptions.ProblemDetails {
	problems := []exceptions.ProblemDetails{
		exceptions.NewProblemDetails(errorType, exceptions.ErrorMessage{
			Title:  title,
			Detail: err.Error(),
		}),
	}

	logger.Log(Logger{
		Context:  ctx,
		Code:     code,
		Message:  title,
		From:     from,
		Layer:    layer,
		TypeLog:  LoggerTypes.ERROR,
		Error:    err,
		Problems: problems,
	})

	return problems
}

func BadRequest(logger LoggerService, ctx context.Context, from, title string, err error) []exceptions.ProblemDetails {
	return LogWithProblem(logger, ctx, from, LoggerLayers.USECASES, title, err, exceptions.RFC400_CODE, exceptions.BadRequest)
}

func Unauthorized(logger LoggerService, ctx context.Context, from, title string, err error) []exceptions.ProblemDetails {
	return LogWithProblem(logger, ctx, from, LoggerLayers.USECASES, title, err, exceptions.RFC401_CODE, exceptions.Unauthorized)
}

func Forbidden(logger LoggerService, ctx context.Context, from, title string, err error) []exceptions.ProblemDetails {
	return LogWithProblem(logger, ctx, from, LoggerLayers.USECASES, title, err, exceptions.RFC403_CODE, exceptions.Forbidden)
}

func NotFound(logger LoggerService, ctx context.Context, from, title string, err error) []exceptions.ProblemDetails {
	return LogWithProblem(logger, ctx, from, LoggerLayers.USECASES, title, err, exceptions.RFC404_CODE, exceptions.NotFound)
}

func Conflict(logger LoggerService, ctx context.Context, from, title string, err error) []exceptions.ProblemDetails {
	return LogWithProblem(logger, ctx, from, LoggerLayers.USECASES, title, err, exceptions.RFC409_CODE, exceptions.Conflict)
}

func InternalServerError(logger LoggerService, ctx context.Context, from, title string, err error) []exceptions.ProblemDetails {
	return LogWithProblem(logger, ctx, from, LoggerLayers.USECASES, title, err, exceptions.RFC500_CODE, exceptions.InternalServerError)
}
