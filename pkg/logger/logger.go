package log

import "go.uber.org/zap"

func Debug(args ...any) {
	zap.S().Debug(args...)
}

func Debugf(template string, args ...any) {
	zap.S().Debugf(template, args...)
}

func Info(args ...any) {
	zap.S().Info(args...)
}

func Infof(template string, args ...any) {
	zap.S().Infof(template, args...)
}

func Warn(args ...any) {
	zap.S().Warn(args...)
}

func Warnf(template string, args ...any) {
	zap.S().Warnf(template, args...)
}

func Error(args ...any) {
	zap.S().Error(args...)
}

func Errorf(template string, args ...any) {
	zap.S().Errorf(template, args...)
}

func Panic(args ...any) {
	zap.S().Panic(args...)
}

func Panicf(template string, args ...any) {
	zap.S().Panicf(template, args...)
}

func Fatal(args ...any) {
	zap.S().Fatal(args...)
}

func Fatalf(template string, args ...any) {
	zap.S().Fatalf(template, args...)
}
