package logging

import (
	"fmt"
	"go-telemetry/pkg/internal/config"
	"sync"
	"time"
)

// A MetaData holds the log variables
type MetaData = map[string]any

// A LoggerData is a user defined log, that takes the timestamp of when the log was initialized
type LoggerData struct {
	LoggerLevel loggerLevel `json:"loggerLevel"`
	Timestamp   time.Time   `json:"timestamp"`
	Message     string      `json:"message"`
	MetaData    MetaData    `json:"metaData"`
}

// A logging holds the top-level configuration of the logger.
// They can be configured via the YAML configuration file or by pre-defined "drivers" or self-created ones.
type logging struct {
	loggerLevel loggerLevel
	outputWrite LogOutputWriter
}

var loggerOnce sync.Once
var loggerInstance *logging
var writeLogOutputMutex sync.Mutex

// NewLog creates a logging instance, respective to the defined YAML configuration or by given "drivers" in form of options argument.
// The "drivers" have a higher priority than YAML configuration.
//
// The log configuration cannot be change once set.
//
// Default:
// If no configuration nor drivers are specified, the log level is Info with output to CLI.
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

// WithLoggerLevel is a pre-defined "driver" that specifies the log level used
func WithLoggerLevel(loggerLevel loggerLevel) func(*logging) {
	return func(l *logging) {
		l.loggerLevel = loggerLevel
	}
}

// WithLogOutputWriter is a pre-defined "driver" that specifies the output writer used
func WithLogOutputWriter(outputWriter LogOutputWriter) func(*logging) {
	return func(l *logging) {
		l.outputWrite = outputWriter
	}
}

// Info registers a log of level Info, with a message and additional attributes.
//
// If no attributes, use nil as MetaData
func (l *logging) Info(msg string, v MetaData) {
	l.processLoggerData(LevelInfo, msg, v)
}

// Warning registers a log of level Warning, with a message and additional attributes.
//
// If no attributes, use nil as MetaData
func (l *logging) Warning(msg string, v MetaData) {
	l.processLoggerData(LevelWarning, msg, v)
}

// Error registers a log of level Error, with a message and additional attributes.
//
// If no attributes, use nil as MetaData
func (l *logging) Error(msg string, v MetaData) {
	l.processLoggerData(LevelError, msg, v)
}

// Debug registers a log of level Debug, with a message and additional attributes.
//
// If no attributes, use nil as MetaData
func (l *logging) Debug(msg string, v MetaData) {
	l.processLoggerData(LevelDebug, msg, v)
}

// processLoggerData initiates the writing of log data to the output by using the specified OutputWriter.
// Logs are printed if the set level is higher or equal than the log method used (Info, Warning, Error, Debug).
// Will block until the writing is finished.
//
// processLoggerData is safe to call concurrently with other operations and will
// block until all other operations finish.
//
// If LogLevel is Off, no logs are printed.
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
