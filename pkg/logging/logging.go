package logging

import (
	"fmt"
	"go-telemetry/config"
	"sync"
	"time"
)

const (
	// Logger levels
	LevelOff     LoggerLevel = "off"     // 0
	LevelInfo    LoggerLevel = "info"    // 1
	LevelWarning LoggerLevel = "warning" // 2
	LevelError   LoggerLevel = "error"   // 3
	LevelDebug   LoggerLevel = "debug"   // 4
)

const (
	// Logger levels
	LevelOffInt     int = 0
	LevelInfoInt    int = 1
	LevelWarningInt int = 2
	LevelErrorInt   int = 3
	LevelDebugInt   int = 4
)

type LoggerLevel string

type MetaData = map[string]any

type LoggerData struct {
	LoggerLevel LoggerLevel
	Timestamp   time.Time
	Message     string
	MetaData    MetaData
}

type Log struct {
	loggerLevel LoggerLevel
	outputWrite OutputWriter
}

var loggerOnce sync.Once
var loggerInstance *Log

func New(options ...func(*Log)) *Log {
	loggerOnce.Do(func() {
		config.Init()

		loggerInstance = &Log{}

		switch config.LoggerConfig.Logger.Level {
		case string(LevelOff), string(LevelInfo), string(LevelWarning), string(LevelError), string(LevelDebug):
			loggerInstance.loggerLevel = LoggerLevel(config.LoggerConfig.Logger.Level)
		default:
			loggerInstance.loggerLevel = LevelInfo
		}

		switch config.LoggerConfig.Logger.OutputWriter {
		case string(cli):
			loggerInstance.outputWrite = CLIOutputWrite()
		case string(jsonFile):
			loggerInstance.outputWrite = JSONOutputFileWrite()
		case string(textFile):
			loggerInstance.outputWrite = TextOutputFileWrite()
		default:
			loggerInstance.outputWrite = CLIOutputWrite()
		}

		// Log options override the YAML file configuration
		for _, o := range options {
			o(loggerInstance)
		}
	})
	return loggerInstance
}

func WithLoggerLevel(loggerLevel LoggerLevel) func(*Log) {
	return func(l *Log) {
		l.loggerLevel = loggerLevel
	}
}

func WithOutputWriter(outputWriter OutputWriter) func(*Log) {
	return func(l *Log) {
		l.outputWrite = outputWriter
	}
}

func convertLoggerLevelToInt(loggerLevel LoggerLevel) int {
	switch loggerLevel {
	case LevelOff:
		return 0
	case LevelInfo:
		return 1
	case LevelWarning:
		return 2
	case LevelError:
		return 3
	case LevelDebug:
		return 4
	default:
		return 1
	}
}

func (l *Log) Info(msg string, v MetaData) {
	l.processLoggerData(LevelInfo, msg, v)
}

func (l *Log) Warning(msg string, v MetaData) {
	l.processLoggerData(LevelWarning, msg, v)
}

func (l *Log) Error(msg string, v MetaData) {
	l.processLoggerData(LevelError, msg, v)
}

func (l *Log) Debug(msg string, v MetaData) {
	l.processLoggerData(LevelDebug, msg, v)
}

func (l *Log) processLoggerData(loggerLevel LoggerLevel, msg string, metaData MetaData) {
	if convertLoggerLevelToInt(loggerLevel) <= convertLoggerLevelToInt(l.loggerLevel) {
		err := l.outputWrite(&LoggerData{
			Timestamp:   time.Now(),
			LoggerLevel: loggerLevel,
			Message:     msg,
			MetaData:    metaData,
		})
		if err != nil {
			fmt.Println(err)
		}
	}
}
