package logger

import (
	"os"
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

type zapLogger struct {
	sugaredLogger *zap.SugaredLogger
}

func getEncoder(isJSON bool) zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time"
	if isJSON {
		return zapcore.NewJSONEncoder(encoderConfig)
	}
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getZapLevel(level string) zapcore.Level {
	switch level {
	case Info:
		return zapcore.InfoLevel
	case Warn:
		return zapcore.WarnLevel
	case Debug:
		return zapcore.DebugLevel
	case Error:
		return zapcore.ErrorLevel
	case Fatal:
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

func newZapLogger(config Configuration) (Logger, error) {
	cores := []zapcore.Core{}

	if config.EnableConsole {
		level := getZapLevel(config.ConsoleLevel)
		writer := zapcore.Lock(os.Stdout)
		core := zapcore.NewCore(getEncoder(config.ConsoleJSONFormat), writer, level)
		cores = append(cores, core)
	}

	if config.EnableFile {
		level := getZapLevel(config.FileLevel)
		writer := zapcore.AddSync(&lumberjack.Logger{
			Filename: config.FileLocation,
			MaxSize:  100,
			Compress: true,
			MaxAge:   28,
		})
		core := zapcore.NewCore(getEncoder(config.FileJSONFormat), writer, level)
		cores = append(cores, core)
	}

	combinedCore := zapcore.NewTee(cores...)

	logger := zap.New(combinedCore,
		zap.AddCallerSkip(2),
		zap.AddCaller(),
	).Sugar()

	return &zapLogger{
		sugaredLogger: logger,
	}, nil
}

func (l *zapLogger) Debug(v ...interface{}) {
	l.sugaredLogger.Debug(l.getMessage(v...))
}

func (l *zapLogger) Info(v ...interface{}) {
	l.sugaredLogger.Info(l.getMessage(v...))
}

func (l *zapLogger) Warn(v ...interface{}) {
	l.sugaredLogger.Warn(l.getMessage(v...))
}

func (l *zapLogger) Error(v ...interface{}) {
	l.sugaredLogger.Error(l.getMessage(v...))
}

func (l *zapLogger) Fatal(v ...interface{}) {
	l.sugaredLogger.Fatal(l.getMessage(v...))
}

func (l *zapLogger) Panic(v ...interface{}) {
	l.sugaredLogger.Fatal(l.getMessage(v...))
}

func (l *zapLogger) WithFields(fields Fields) Logger {
	var f = make([]interface{}, 0)
	for k, v := range fields {
		f = append(f, k)
		f = append(f, v)
	}
	newLogger := l.sugaredLogger.With(f...)
	return &zapLogger{newLogger}
}

func (l *zapLogger) getMessage(v ...interface{}) (sentence string) {
	for i := range v {
		var msg string
		if str, ok := v[i].(fmt.Stringer); ok {
			msg = str.String()
		} else {
			switch t := v[i].(type) {
			case string:
				msg = t
			case error:
				msg = t.Error()
			default:
				// TODO
				msg = fmt.Sprint(v[i])
			}
		}

		if sentence != "" {
			sentence += msg
		} else {
			sentence = msg
		}
	}

	return sentence
}
