package logging

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"pet-services-api/internal/exceptions"

	"github.com/lmittmann/tint"
)

const TimeFormat = "2006-01-02 15:04:05"

type DefaultLogger struct{}

type LoggerInterface interface {
	LogError(ctx context.Context, from, message string, err error)
	LogInfo(ctx context.Context, from, message string)
	LogWarning(ctx context.Context, from, message string, err error)
	LogBadRequest(ctx context.Context, from, title string, err error) []exceptions.ProblemDetails
	LogMultipleBadRequests(ctx context.Context, from, title string, problems []exceptions.ProblemDetails)
	LogUnauthorized(ctx context.Context, from, title string, err error) []exceptions.ProblemDetails
	LogForbidden(ctx context.Context, from, title string, err error) []exceptions.ProblemDetails
	LogNotFound(ctx context.Context, from, title string, err error) []exceptions.ProblemDetails
	LogConflict(ctx context.Context, from, title string, err error) []exceptions.ProblemDetails
	LogInternalServerError(ctx context.Context, from, title string, err error) []exceptions.ProblemDetails
}

func (l *DefaultLogger) LogMultipleBadRequests(ctx context.Context, from, title string, problems []exceptions.ProblemDetails) {
	for _, p := range problems {
		l.LogError(ctx, from, title, errors.New(p.Detail))
	}
}

func (l *DefaultLogger) LogBadRequest(ctx context.Context, from, title string, err error) []exceptions.ProblemDetails {
	problems := []exceptions.ProblemDetails{
		exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  title,
			Detail: err.Error(),
		}),
	}
	l.LogError(ctx, from, title, err)
	return problems
}

func (l *DefaultLogger) LogUnauthorized(ctx context.Context, from, title string, err error) []exceptions.ProblemDetails {
	problems := []exceptions.ProblemDetails{
		exceptions.NewProblemDetails(exceptions.Unauthorized, exceptions.ErrorMessage{
			Title:  title,
			Detail: err.Error(),
		}),
	}
	l.LogError(ctx, from, title, err)
	return problems
}

func (l *DefaultLogger) LogForbidden(ctx context.Context, from, title string, err error) []exceptions.ProblemDetails {
	problems := []exceptions.ProblemDetails{
		exceptions.NewProblemDetails(exceptions.Forbidden, exceptions.ErrorMessage{
			Title:  title,
			Detail: err.Error(),
		}),
	}
	l.LogError(ctx, from, title, err)
	return problems
}

func (l *DefaultLogger) LogNotFound(ctx context.Context, from, title string, err error) []exceptions.ProblemDetails {
	problems := []exceptions.ProblemDetails{
		exceptions.NewProblemDetails(exceptions.NotFound, exceptions.ErrorMessage{
			Title:  title,
			Detail: err.Error(),
		}),
	}
	l.LogError(ctx, from, title, err)
	return problems
}

func (l *DefaultLogger) LogConflict(ctx context.Context, from, title string, err error) []exceptions.ProblemDetails {
	problems := []exceptions.ProblemDetails{
		exceptions.NewProblemDetails(exceptions.Conflict, exceptions.ErrorMessage{
			Title:  title,
			Detail: err.Error(),
		}),
	}
	l.LogError(ctx, from, title, err)
	return problems
}

func (l *DefaultLogger) LogInternalServerError(ctx context.Context, from, title string, err error) []exceptions.ProblemDetails {
	problems := []exceptions.ProblemDetails{
		exceptions.NewProblemDetails(exceptions.InternalServerError, exceptions.ErrorMessage{
			Title:  title,
			Detail: err.Error(),
		}),
	}
	l.LogError(ctx, from, title, err)
	return problems
}

func (l *DefaultLogger) LogError(ctx context.Context, from, message string, err error) {
	if logger == nil {
		InitLogger()
	}
	logger.ErrorContext(
		ctx,
		"ERROR",
		"message:", message,
		"from:", from,
		"error:", err,
	)
}

func (l *DefaultLogger) LogInfo(ctx context.Context, from, message string) {
	if logger == nil {
		InitLogger()
	}
	logger.InfoContext(
		ctx,
		"INFO",
		"message:", message,
		"from:", from,
	)
}

func (l *DefaultLogger) LogWarning(ctx context.Context, from, message string, err error) {
	if logger == nil {
		InitLogger()
	}
	logger.WarnContext(
		ctx,
		"WARNING",
		"message:", message,
		"from:", from,
		"error:", err,
	)
}

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

var logger *slog.Logger

func InitLogger() {
	logger = slog.New(tint.NewHandler(os.Stdout, &tint.Options{
		Level:      slog.LevelDebug,
		TimeFormat: TimeFormat,
	}))
}

func NewLogger(log Logger) {
	if logger == nil {
		InitLogger()
	}

	switch log.TypeLog {
	case "ERROR":
		logger.ErrorContext(
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
		logger.InfoContext(
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
		logger.WarnContext(
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
