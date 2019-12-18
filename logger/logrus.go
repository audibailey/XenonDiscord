package logger

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

type logrusLogEntry struct {
	entry *logrus.Entry
}

type logrusLogger struct {
	logger *logrus.Logger
}

func getFormatter(isJSON bool) logrus.Formatter {
	if isJSON {
		return &logrus.JSONFormatter{}
	}
	return &logrus.TextFormatter{
		FullTimestamp:          true,
		DisableLevelTruncation: true,
	}
}

func newLogrusLogger(config Configuration) (Logger, error) {
	logLevel := config.ConsoleLevel
	if logLevel == "" {
		logLevel = config.FileLevel
	}

	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		return nil, err
	}

	stdOutHandler := os.Stdout
	fileHandler := &lumberjack.Logger{
		Filename: config.FileLocation,
		MaxSize:  100,
		Compress: true,
		MaxAge:   28,
	}
	lLogger := &logrus.Logger{
		Out:       stdOutHandler,
		Formatter: getFormatter(config.ConsoleJSONFormat),
		Hooks:     make(logrus.LevelHooks),
		Level:     level,
	}

	if config.EnableConsole && config.EnableFile {
		lLogger.SetOutput(io.MultiWriter(stdOutHandler, fileHandler))
	} else {
		if config.EnableFile {
			lLogger.SetOutput(fileHandler)
			lLogger.SetFormatter(getFormatter(config.FileJSONFormat))
		}
	}

	return &logrusLogger{
		logger: lLogger,
	}, nil
}

func (l *logrusLogger) Debug(v ...interface{}) {
	l.logger.Debug(v...)
}

func (l *logrusLogger) Info(v ...interface{}) {
	l.logger.Info(v...)
}

func (l *logrusLogger) Warn(v ...interface{}) {
	l.logger.Warn(v...)
}

func (l *logrusLogger) Error(v ...interface{}) {
	l.logger.Error(v...)
}

func (l *logrusLogger) Fatal(v ...interface{}) {
	l.logger.Fatal(v...)
}

func (l *logrusLogger) Panic(v ...interface{}) {
	l.logger.Fatal(v)
}

func (l *logrusLogger) WithFields(fields Fields) Logger {
	return &logrusLogEntry{
		entry: l.logger.WithFields(convertToLogrusFields(fields)),
	}
}

func (l *logrusLogEntry) Debug(v ...interface{}) {
	l.Debug(v...)
}

func (l *logrusLogEntry) Info(v ...interface{}) {
	l.Info(v...)
}

func (l *logrusLogEntry) Warn(v ...interface{}) {
	l.Warn(v...)
}

func (l *logrusLogEntry) Error(v ...interface{}) {
	l.Error(v...)
}

func (l *logrusLogEntry) Fatal(v ...interface{}) {
	l.Fatal(v...)
}

func (l *logrusLogEntry) Panic(v ...interface{}) {
	l.Fatal(v...)
}

func (l *logrusLogEntry) WithFields(fields Fields) Logger {
	return l.WithFields(fields)
}

func convertToLogrusFields(fields Fields) logrus.Fields {
	logrusFields := logrus.Fields{}
	for index, val := range fields {
		logrusFields[index] = val
	}
	return logrusFields
}
