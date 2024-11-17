package logging

import (
	"fmt"
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

const (
	cli      OutputWriterType = "cli"
	jsonFile OutputWriterType = "jsonFile"
	textFile OutputWriterType = "textFile"
)

type LoggerLevel string

type OutputWriterType string

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
