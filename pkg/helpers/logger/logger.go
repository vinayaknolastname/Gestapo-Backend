package logger

type Logger interface {
	LogInfo(args ...interface{})
	LogError(args ...interface{})
	LogPanic(args ...interface{})
	LogFatal(args ...interface{})
}
