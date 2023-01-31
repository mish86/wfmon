package log

import (
	"fmt"
	"os"
	"time"
	app "wfmon/pkg"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	callerSkipLevel = 1

	logFilename   = "/usr/local/var/log/wfmon.log"
	logMaxSize    = 500 // megabytes
	logMaxBackups = 3
	logMaxAge     = 28 // days
)

// Creates logger with config depending on application mode: dev / prod.
func NewLogger(mode app.Mode) *zap.Logger {
	core, options := getCore(mode)
	logger := zap.New(core, options...)
	// logger := zap.NewNop()

	return logger
}

func NewSLogger(logger *zap.Logger) *zap.SugaredLogger {
	return logger.Sugar()
}

func getDevCore() zapcore.Core {
	// level-handling logic
	highPriority := zap.NewAtomicLevelAt(zap.ErrorLevel)
	lowPriority := zap.NewAtomicLevelAt(zap.DebugLevel)

	// Encoder Configuration
	encoderCfg := zap.NewDevelopmentEncoderConfig()
	encoderCfg.EncodeTime = syslogTimeEncoder
	encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoderCfg.EncodeCaller = zapcore.ShortCallerEncoder
	encoderCfg.ConsoleSeparator = " "

	// syslog output for human operators.
	syslogEncoder := zapcore.NewConsoleEncoder(encoderCfg)

	// High-priority output should also go to standard error, and low-priority
	// output should also go to standard out.
	stderr := zapcore.Lock(os.Stderr)
	writer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logFilename,
		MaxSize:    logMaxSize,
		MaxBackups: logMaxBackups,
		MaxAge:     logMaxAge,
	})

	core := zapcore.NewTee(
		zapcore.NewCore(syslogEncoder, stderr, highPriority),
		zapcore.NewCore(syslogEncoder, writer, lowPriority),
	)

	return core
}
func getProdCore() zapcore.Core {
	// level-handling logic
	highPriority := zap.NewAtomicLevelAt(zap.ErrorLevel)
	lowPriority := zap.NewAtomicLevelAt(zap.DebugLevel)

	// Encoder Configuration
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = syslogTimeEncoder
	encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoderCfg.EncodeCaller = nil
	encoderCfg.ConsoleSeparator = " "

	// syslog output for human operators.
	syslogEncoder := zapcore.NewConsoleEncoder(encoderCfg)

	// High-priority output should also go to standard error, and low-priority
	// output should also go to standard out.
	stderr := zapcore.Lock(os.Stderr)
	// consoleDebugging := zapcore.Lock(os.Stdout)
	// consoleDebugging := zapcore.Lock(os.Stdout)
	writer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logFilename,
		MaxSize:    logMaxSize,
		MaxBackups: logMaxBackups,
		MaxAge:     logMaxAge,
	})

	core := zapcore.NewTee(
		zapcore.NewCore(syslogEncoder, stderr, highPriority),
		zapcore.NewCore(syslogEncoder, writer, lowPriority),
	)

	return core
}

func getCore(mode app.Mode) (zapcore.Core, []zap.Option) {
	switch mode {
	case app.Dev:
		return getDevCore(), []zap.Option{
			zap.AddCaller(),
			zap.AddCallerSkip(callerSkipLevel),
		}
	case app.Prod:
		return getProdCore(), []zap.Option{}
	default:
		panic(fmt.Errorf("got unsupported application mode %s", mode))
	}
}

func syslogTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(time.Stamp))
}
