package logging

import (
	"fmt"
	"go-telemetry/config"
	"sync"
	"time"
)

type MetaData = map[string]any

type LoggerData struct {
	LoggerLevel LoggerLevel `json:"loggerLevel"`
	Timestamp   time.Time   `json:"timestamp"`
	Message     string      `json:"message"`
	MetaData    MetaData    `json:"metaData"`
}

type Log struct {
	loggerLevel LoggerLevel
	outputWrite LogOutputWriter
}

var loggerOnce sync.Once
var loggerInstance *Log

func NewLog(options ...func(*Log)) *Log {
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
			loggerInstance.outputWrite = CLILogOutputWrite()
		case string(jsonFile):
			loggerInstance.outputWrite = JSONLogOutputFileWrite()
		case string(textFile):
			loggerInstance.outputWrite = TextLogOutputFileWrite()
		default:
			loggerInstance.outputWrite = CLILogOutputWrite()
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

func WithLogOutputWriter(outputWriter LogOutputWriter) func(*Log) {
	return func(l *Log) {
		l.outputWrite = outputWriter
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
