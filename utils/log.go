package utils

import (
	"context"
	"log"

	// "github.com/kostyay/zapdriver"
	"github.com/blendle/zapdriver"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

const (
	traceKey        = "logging.googleapis.com/trace"
	spanKey         = "logging.googleapis.com/spanId"
	traceSampledKey = "logging.googleapis.com/trace_sampled"
	errorKey        = "err"
)

var googleProjectID = "tilda-trial-dev"

type Logger struct {
	logger *zap.SugaredLogger
}

var globalLog *Logger

func LoggerInit() {
	zapconf := zapdriver.NewProductionConfig()

	l, err := zapconf.Build(
		zapdriver.WrapCore(
			zapdriver.ReportAllErrors(true),
			zapdriver.ServiceName(ServiceName),
			// zapdriver.ServiceVersion(Version),
		),
		zap.AddCallerSkip(1),
	)
	if err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}

	globalLog = &Logger{logger: l.Sugar()}
}

func (l *Logger) WithContext(ctx context.Context) *Logger {
	span := trace.SpanFromContext(ctx).SpanContext()

	return &Logger{
		logger: l.logger.With(
			traceKey, traceID(span.TraceID().String()),
			spanKey, span.SpanID().String(),
			traceSampledKey, span.IsSampled(),
		)}
}

func traceID(id string) string {
	return "projects/" + googleProjectID + "/traces/" + id
}

func (l *Logger) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

func (l *Logger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.logger.Error(args...)
}

func (l *Logger) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}

func (l *Logger) WithError(err error) *Logger {
	return &Logger{logger: l.logger.With(errorKey, err.Error())}
}

func WithContext(ctx context.Context) *Logger {
	return globalLog.WithContext(ctx)
}

func Debug(args ...interface{}) {
	globalLog.Debug(args...)
}

func Info(args ...interface{}) {
	globalLog.Info(args...)
}

func Infof(format string, args ...interface{}) {
	globalLog.logger.Infof(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	globalLog.logger.Fatalf(format, args...)
}

func Fatal(args ...interface{}) {
	globalLog.logger.Fatal(args...)
}

func WithError(err error) *Logger {
	return &Logger{logger: globalLog.logger.With(errorKey, err.Error())}
}
