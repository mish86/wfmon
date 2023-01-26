package log

import (
	"fmt"
	"os"
	"time"
	"wfmon/pkg"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const callerSkipLevel = 1

// Creates logger with config depending on application mode: dev / prod.
func NewLogger(mode pkg.Mode) *zap.Logger {
	core, options := getCore(mode)
	logger := zap.New(core, options...)

	return logger
}

func NewSLogger(logger *zap.Logger) *zap.SugaredLogger {
	return logger.Sugar()
}

func getDevCore() zapcore.Core {
	// level-handling logic
	lowPriority := zap.NewAtomicLevelAt(zap.DebugLevel)

	// Encoder Configuration
	encoderCfg := zap.NewDevelopmentEncoderConfig()
	// encoderCfg.EncodeTime = zapcore.RFC3339TimeEncoder
	encoderCfg.EncodeTime = syslogTimeEncoder
	encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoderCfg.EncodeCaller = zapcore.ShortCallerEncoder
	encoderCfg.ConsoleSeparator = " "

	// console output for human operators.
	consoleEncoder := zapcore.NewConsoleEncoder(encoderCfg)

	// High-priority output should also go to standard error, and low-priority
	// output should also go to standard out.
	consoleDebugging := zapcore.Lock(os.Stdout)

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, consoleDebugging, lowPriority),
	)

	return core
}

func getProdCore() zapcore.Core {
	// level-handling logic
	highPriority := zap.NewAtomicLevelAt(zap.ErrorLevel)
	lowPriority := zap.NewAtomicLevelAt(zap.DebugLevel)

	// Encoder Configuration
	encoderCfg := zap.NewProductionEncoderConfig()
	// encoderCfg.EncodeTime = zapcore.RFC3339TimeEncoder
	encoderCfg.EncodeTime = syslogTimeEncoder
	encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoderCfg.EncodeCaller = zapcore.ShortCallerEncoder
	encoderCfg.ConsoleSeparator = " "

	// console output for human operators.
	consoleEncoder := zapcore.NewConsoleEncoder(encoderCfg)

	// High-priority output should also go to standard error, and low-priority
	// output should also go to standard out.
	consoleDebugging := zapcore.Lock(os.Stdout)
	consoleErrors := zapcore.Lock(os.Stderr)

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, consoleErrors, highPriority),
		zapcore.NewCore(consoleEncoder, consoleDebugging, lowPriority),
	)

	return core
}

func getCore(mode pkg.Mode) (zapcore.Core, []zap.Option) {
	switch mode {
	case pkg.Dev:
		return getDevCore(), []zap.Option{
			zap.AddCaller(),
			zap.AddCallerSkip(callerSkipLevel),
		}
	case pkg.Prod:
		return getProdCore(), []zap.Option{}
	default:
		panic(fmt.Errorf("got unsupported application mode %s", mode))
	}
}

func syslogTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(time.Stamp))
}
