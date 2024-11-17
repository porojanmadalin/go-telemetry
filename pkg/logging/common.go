package logging

const (
	// Logger levels
	LevelOff     loggerLevel = "off"     // 0
	LevelInfo    loggerLevel = "info"    // 1
	LevelWarning loggerLevel = "warning" // 2
	LevelError   loggerLevel = "error"   // 3
	LevelDebug   loggerLevel = "debug"   // 4
)

const (
	// Logger levels
	levelOffInt     int = 0
	levelInfoInt    int = 1
	levelWarningInt int = 2
	levelErrorInt   int = 3
	levelDebugInt   int = 4
)

const (
	cli      OutputWriterType = "cli"
	jsonFile OutputWriterType = "jsonFile"
	textFile OutputWriterType = "textFile"
)

type loggerLevel string

type OutputWriterType string

func convertLoggerLevelToInt(loggerLevel loggerLevel) int {
	switch loggerLevel {
	case LevelOff:
		return levelOffInt
	case LevelInfo:
		return levelInfoInt
	case LevelWarning:
		return levelWarningInt
	case LevelError:
		return levelErrorInt
	case LevelDebug:
		return levelDebugInt
	default:
		return levelOffInt
	}
}
