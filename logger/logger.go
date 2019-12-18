package logger

import (
	"errors"
	glog "log"

	"github.com/audibailey/XenonDiscord/config"

	"github.com/spf13/viper"
	"github.com/mitchellh/go-homedir"
)

var Log Logger = (*Empty)(nil)

type Fields map[string]interface{}

const (
	Debug = "debug"
	Info = "info"
	Warn = "warn"
	Error = "error"
	Fatal = "fatal"
)

const (
	InstanceZapLogger int = iota
	InstanceLogrusLogger int = iota
)

var (
	errInvalidLoggerInstance = errors.New("Invalid logger instance")
)

type Empty struct{}


type Logger interface {
	Debug(v ...interface{})

	Info(v ...interface{})

	Warn(v ...interface{})

	Error(v ...interface{})

	Fatal(v ...interface{})

	Panic(v ...interface{})	

	WithFields(keyValues Fields) Logger
}

type Configuration struct {
	EnableConsole     bool
	ConsoleJSONFormat bool
	ConsoleLevel      string
	EnableFile        bool
	FileJSONFormat    bool
	FileLevel         string
	FileLocation      string
}


func NewLogger(config Configuration, loggerInstance int) error {
	if loggerInstance == InstanceZapLogger {
		logger, err := newZapLogger(config)
		if err != nil {
			return err
		}
		Log = logger
		return nil
	} else if loggerInstance == InstanceLogrusLogger {
		logger, err := newLogrusLogger(config)
		if err != nil {
			return err
		}
		Log = logger
		return nil
	}
	return errInvalidLoggerInstance
}

func (Empty) Debug(v ...interface{}) {
	if viper.GetBool("debug") {
		Log.Debug(v...)
	}
}

func (Empty) Info(v ...interface{}) {
	Log.Info(v...)
}

func (Empty) Warn(v ...interface{}) {
	Log.Warn(v...)
}

func (Empty) Error(v ...interface{}) {
	Log.Error(v...)
}

func (Empty) Fatal(v ...interface{}) {
	Log.Fatal(v)
}

func (Empty) Panic(v ...interface{}) {
	Log.Panic(v)
}

func (Empty) WithFields(keyValues Fields) Logger {
	return Log.WithFields(keyValues)
}

func Configure() (Logger) {

	fileloc, err := homedir.Expand(config.Conf.Logging.FileLocation)
	if err != nil {
		glog.Fatal("Could not find file")
	}

	logConfig := Configuration{
		EnableConsole:     config.Conf.Logging.EnableConsole,
		ConsoleLevel:      config.Conf.Logging.ConsoleLevel,
		ConsoleJSONFormat: config.Conf.Logging.ConsoleJSONFormat,
		EnableFile:        config.Conf.Logging.EnableFile,
		FileLevel:         config.Conf.Logging.FileLevel,
		FileJSONFormat:    config.Conf.Logging.FileJSONFormat,
		FileLocation:      fileloc,
	}

	err = NewLogger(logConfig, config.Conf.Logging.Logger)
	if err != nil {
		glog.Fatal("Could not instantiate log %s", err.Error())
	}

	return Log
}
