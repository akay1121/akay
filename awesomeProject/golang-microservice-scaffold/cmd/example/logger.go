package main

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"

	kratoszap "github.com/go-kratos/kratos/contrib/log/zap/v2"
)

// NewLogger creates a logger based on the configuration
//
// Note that the logger should typically write multiple logs simultaneously to different
// outputs (Standard Output [os.Stdout], Log files, or Log collectors like logstash, filebeat, etc.).
func NewLogger(driver, addr string) log.Logger {
	return log.With(
		kratoszap.NewLogger(
			zap.New(
				zapcore.NewTee(
					zapcore.NewCore(
						zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
							CallerKey:        "caller",
							LevelKey:         "level",
							MessageKey:       "msg",
							TimeKey:          "ts",
							StacktraceKey:    "st",
							LineEnding:       zapcore.DefaultLineEnding,
							ConsoleSeparator: " ",
							EncodeCaller: func(caller zapcore.EntryCaller, encoder zapcore.PrimitiveArrayEncoder) {
								encoder.AppendString(
									fmt.Sprintf(
										"\033[36m%s:%d\033[0m",
										caller.File,
										caller.Line,
									),
								)
							},
							// Print out the timestamp with customized formatting
							EncodeTime: func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
								encoder.AppendString(t.Format(timeFormat))
							},
							// Print out different levels using different colors
							EncodeLevel: func(level zapcore.Level, encoder zapcore.PrimitiveArrayEncoder) {
								switch level {
								case zapcore.DebugLevel:
									encoder.AppendString("\033[37;1mDEBUG\033[0m")
								case zapcore.InfoLevel:
									encoder.AppendString("\033[36;1mINFO\033[0m")
								case zapcore.WarnLevel:
									encoder.AppendString("\033[33;1mWARN\033[0m")
								case zapcore.ErrorLevel:
									encoder.AppendString("\033[31;1mERROR\033[0m")
								case zapcore.DPanicLevel:
									encoder.AppendString("\033[31;1mDPANIC\033[0m")
								case zapcore.PanicLevel:
									encoder.AppendString("\033[31;1mPANIC\033[0m")
								case zapcore.FatalLevel:
									encoder.AppendString("\033[31;1mFATAL\033[0m")
								case zapcore.InvalidLevel:
									encoder.AppendString("\033[31;1mINVALID\033[0m")
								}
								encoder.AppendString(fmt.Sprintf("\033[35m[%s]\033[0m", id))
							},
						}),
						// We need to tell the zap core that we want the log message to be printed to the console
						zapcore.AddSync(os.Stdout),
						// All the log messages whose level is beyond [zapcore.DebugLevel] are printed
						zapcore.DebugLevel,
					),
					// You can add one more output path as well by uncommenting the following lines.
					//zapcore.NewCore(
					//	zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
					//	zapcore.AddSync(logOutput),
					//	zapcore.DebugLevel,
					//),
				),
				// Print out the exact location of the caller
				zap.WithCaller(true),
				// Print out the stack trace when the level is greater than or equal to [zap.WarnLevel]
				zap.AddStacktrace(zap.WarnLevel),
				zap.Fields(
					zap.String("service.id", id),
					zap.String("service.name", Name),
					zap.String("service.version", Version),
				),
			),
		),
		// Each call to the service should have its tracing id, so we also need print it out.
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)
}
