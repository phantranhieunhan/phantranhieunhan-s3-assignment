package logger

import (
	"github.com/phantranhieunhan/s3-assignment/common/constant"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func init() {
	l = zap.NewNop().Sugar()
}

// Setup: creates a new global logger instance, should be call first when application start.
func Setup(environment string) {
	conf := newProductionConfig()
	if environment != constant.PRODUCTION_ENV_NAME {
		conf = zap.NewDevelopmentConfig()
	}

	conf.DisableStacktrace = true
	log, err := conf.Build()
	if err != nil {
		panic(err)
	}

	l = log.WithOptions(zap.AddCallerSkip(1)).Sugar()
}

// newProductionConfig is a reasonable production logging configuration.
// Logging is enabled at InfoLevel and above.
//
// It uses a CONSOLE encoder, writes to standard error, and enables sampling.
// Stacktraces are automatically included on logs of ErrorLevel and above.
func newProductionConfig() zap.Config {
	return zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "console",
		EncoderConfig:    newProductionEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}
}

// newProductionEncoderConfig returns an opinionated EncoderConfig for
// production environments.
func newProductionEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

// Debug uses fmt.Sprint to construct and log a message
func Debug(args ...interface{}) {
	l.Debug(args...)
}

// Debugf uses fmt.Sprintf to log a templated message
func Debugf(template string, args ...interface{}) {
	l.Debugf(template, args...)
}

// Debugw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With
func Debugw(msg string, keysValues ...interface{}) {
	l.Debugw(msg, keysValues...)
}

// Info uses fmt.Sprint to construct and log a message
func Info(args ...interface{}) {
	l.Info(args...)
}

// Infof uses fmt.Sprintf to log a templated message
func Infof(template string, args ...interface{}) {
	l.Infof(template, args...)
}

// Infow logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func Infow(msg string, keysValues ...interface{}) {
	l.Infow(msg, keysValues...)
}

// Warn uses fmt.Sprint to construct and log a message
func Warn(args ...interface{}) {
	l.Warn(args...)
}

// Warnf uses fmt.Sprintf to log a templated message
func Warnf(template string, args ...interface{}) {
	l.Warnf(template, args...)
}

// Warnw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func Warnw(msg string, keysValues ...interface{}) {
	l.Warnw(msg, keysValues...)
}

// Error uses fmt.Sprint to construct and log a message
func Error(args ...interface{}) {
	l.Error(args...)
}

// Errorf uses fmt.Sprintf to log a templated message
func Errorf(template string, args ...interface{}) {
	l.Errorf(template, args...)
}

// Errorw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func Errorw(msg string, keysValues ...interface{}) {
	l.Errorw(msg, keysValues...)
}

// Fatal uses fmt.Sprint to construct and log a message, then calls os.Exit
func Fatal(args ...interface{}) {
	l.Fatal(args...)
}

// Fatalf uses fmt.Sprintf to log a templated message, then calls os.Exit
func Fatalf(template string, args ...interface{}) {
	l.Fatalf(template, args...)
}

// Fatalw logs a message with some additional context, then calls os.Exit. The
// variadic key-value pairs are treated as they are in With
func Fatalw(msg string, keysValues ...interface{}) {
	l.Fatalw(msg, keysValues...)
}

// Panic uses fmt.Sprint to construct and log a message, then panics
func Panic(args ...interface{}) {
	l.Panic(args...)
}

// Panicf uses fmt.Sprintf to log a templated message, then panics
func Panicf(template string, args ...interface{}) {
	l.Panicf(template, args...)
}

// Panicw logs a message with some additional context, then panics. The
// variadic key-value pairs are treated as they are in With
func Panicw(msg string, keysValues ...interface{}) {
	l.Panicw(msg, keysValues...)
}

// WithLogger set global logger by new logger
func WithLogger(_logger Log) {
	l = _logger
}
