package logging

import (
	"fmt"
	"go-telemetry/pkg/internal/config"
	"sync"
	"time"
)

type MetaData = map[string]any

type LoggerData struct {
	LoggerLevel loggerLevel `json:"loggerLevel"`
	Timestamp   time.Time   `json:"timestamp"`
	Message     string      `json:"message"`
	MetaData    MetaData    `json:"metaData"`
}

type logging struct {
	loggerLevel loggerLevel
	outputWrite LogOutputWriter
}

var loggerOnce sync.Once
var loggerInstance *logging
var writeLogOutputMutex sync.Mutex

func NewLog(options ...func(*logging)) *logging {
	loggerOnce.Do(func() {
		config.Init()

		loggerInstance = &logging{}

		switch config.LoggerConfig.Logger.Level {
		case string(LevelOff), string(LevelInfo), string(LevelWarning), string(LevelError), string(LevelDebug):
			loggerInstance.loggerLevel = loggerLevel(config.LoggerConfig.Logger.Level)
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

func WithLoggerLevel(loggerLevel loggerLevel) func(*logging) {
	return func(l *logging) {
		l.loggerLevel = loggerLevel
	}
}

func WithLogOutputWriter(outputWriter LogOutputWriter) func(*logging) {
	return func(l *logging) {
		l.outputWrite = outputWriter
	}
}

func (l *logging) Info(msg string, v MetaData) {
	l.processLoggerData(LevelInfo, msg, v)
}

func (l *logging) Warning(msg string, v MetaData) {
	l.processLoggerData(LevelWarning, msg, v)
}

func (l *logging) Error(msg string, v MetaData) {
	l.processLoggerData(LevelError, msg, v)
}

func (l *logging) Debug(msg string, v MetaData) {
	l.processLoggerData(LevelDebug, msg, v)
}

func (l *logging) processLoggerData(loggerLevel loggerLevel, msg string, metaData MetaData) {
	if convertLoggerLevelToInt(loggerLevel) <= convertLoggerLevelToInt(l.loggerLevel) {
		writeLogOutputMutex.Lock()
		err := l.outputWrite(&LoggerData{
			Timestamp:   time.Now(),
			LoggerLevel: loggerLevel,
			Message:     msg,
			MetaData:    metaData,
		})
		if err != nil {
			fmt.Println(err)
		}
		writeLogOutputMutex.Unlock()
	}
}
