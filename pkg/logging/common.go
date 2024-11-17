package logging

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

type Loggers struct {
	Log
	TransactionLog
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
