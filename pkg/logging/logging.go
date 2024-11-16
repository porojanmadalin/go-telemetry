package logging

type LoggingLevel int

const (
	// Logging levels
	Off     LoggingLevel = 0
	Info    LoggingLevel = 1
	Warning LoggingLevel = 2
	Error   LoggingLevel = 3
	Debug   LoggingLevel = 4
)

type Logger struct{}
