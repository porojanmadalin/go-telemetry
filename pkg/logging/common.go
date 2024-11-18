package logging

// Logger levels used by logger driver.
const (
	LevelOff     loggerLevel = "off"     // 0
	LevelInfo    loggerLevel = "info"    // 1
	LevelWarning loggerLevel = "warning" // 2
	LevelError   loggerLevel = "error"   // 3
	LevelDebug   loggerLevel = "debug"   // 4
)

// Logger levels numbers used when comparing string log levels
const (
	levelOffInt     int = 0
	levelInfoInt    int = 1
	levelWarningInt int = 2
	levelErrorInt   int = 3
	levelDebugInt   int = 4
)

// Output writers types that resemble the way the logs are printed
const (
	cli      OutputWriterType = "cli"
	jsonFile OutputWriterType = "jsonFile"
	textFile OutputWriterType = "textFile"
)

const (
	// JSON File Writer Constants
	startArray      = "[\n"
	endArray        = "\n]"
	indent          = "  "
	objectDelimiter = ",\n" + indent
)

const (
	// timestamp format
	timestampFormat     = "2006-01-02 15:04:05.0000"
	fileTimestampFormat = "2006-01-02"
)

type loggerLevel string

// A OutputWriterType is a output writer driver identifier.
type OutputWriterType string

// convertLoggerLevelToInt converts loggerLevel type to its corresponding number value
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
		return levelInfoInt
	}
}
