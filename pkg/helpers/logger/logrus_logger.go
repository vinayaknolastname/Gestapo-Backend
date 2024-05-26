package logger

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
	lj "gopkg.in/natefinch/lumberjack.v2"
)

// Log LogrusLogger
type LogrusLogger struct {
	infoLogger  *logrus.Logger
	errorLogger *logrus.Logger
	panicLogger *logrus.Logger
	fatalLogger *logrus.Logger
}

// event stores messages to be logged
type event struct {
	name   string
	level  logrus.Level
	logger *logrus.Logger
	args   interface{}
}

/*
Creates a new Logrus Logger pointer to be used for logging.
This function will need a log ile name.
The log file name should represent the service name.
*/
func NewLogrusLogger(logFileName string) Logger {
	return &LogrusLogger{
		infoLogger:  createLogger("info_"+logFileName, false),
		errorLogger: createLogger("error_"+logFileName, false),
		panicLogger: createLogger("panic_"+logFileName, false),
		fatalLogger: createLogger("fatal_"+logFileName, false),
	}
}

func createLogger(logFileName string, setReportCaller bool) *logrus.Logger {
	var logger = logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetReportCaller(setReportCaller)
	var loggerRotate = &lj.Logger{
		Filename:  getLogPath(logFileName),
		MaxSize:   100,
		MaxAge:    30,
		Compress:  true,
		LocalTime: false,
	}
	mw := io.MultiWriter(os.Stdout, loggerRotate)
	logger.SetOutput(mw)
	return logger
}

// This set's the log path as the current directory where the execution is done from.
func getLogPath(logFileName string) string {
	path, err := os.Getwd()
	if err != nil {
		logrus.Error("Failed to get log file path")
	}
	var logDirPath = path + string(os.PathSeparator) + "logs" + string(os.PathSeparator)
	err = os.MkdirAll(logDirPath, 0755)
	if err != nil {
		logrus.Error("Failed to create dir for logs")
	}
	return logDirPath + logFileName
}

// LogInfo implements Logger.
func (l *LogrusLogger) LogInfo(args ...interface{}) {
	e := &event{
		name:   "Info",
		level:  logrus.InfoLevel,
		logger: l.infoLogger,
		args:   args,
	}
	l.log(e)
}

// LogError implements Logger.
func (l *LogrusLogger) LogError(args ...interface{}) {
	e := &event{
		name:   "Error",
		level:  logrus.ErrorLevel,
		logger: l.errorLogger,
		args:   args,
	}
	l.log(e)
}

// LogFatal implements Logger.
func (l *LogrusLogger) LogFatal(args ...interface{}) {
	event := &event{
		name:   "Fatal",
		level:  logrus.FatalLevel,
		logger: l.fatalLogger,
		args:   args,
	}
	event.logger.Fatal(event.level, event.args)
}

// LogPanic implements Logger.
func (l *LogrusLogger) LogPanic(args ...interface{}) {
	event := &event{
		name:   "Panic",
		level:  logrus.PanicLevel,
		logger: l.panicLogger,
		args:   args,
	}
	event.logger.Panic(event.level, event.args)
}

func (l *LogrusLogger) log(event *event) {
	event.logger.Log(event.level, event.args)
	event = nil
}
