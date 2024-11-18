package logging

import (
	"fmt"
	"go-telemetry/pkg/internal/config"
	itesting "go-telemetry/pkg/internal/telemetrytesting"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	// output writers fn names
	CLILogOutputWriteName      = "CLILogOutputWrite"
	JSONLogOutputFileWriteName = "JSONLogOutputFileWrite"
	TextLogOutputFileWriteName = "TextLogOutputFileWrite"
)

func TestNewLogWithYAMLConfig(t *testing.T) {
	type Expected struct {
		Level            loggerLevel
		OutputWriterName string
	}
	type TestCase struct {
		Data     config.Logger
		Expected Expected
	}

	testCases := []TestCase{
		{
			Data: config.Logger{
				Level:        "undefined",
				OutputWriter: "undefined",
			},
			Expected: Expected{
				Level:            LevelOff,
				OutputWriterName: CLILogOutputWriteName,
			},
		},
		{
			Data: config.Logger{
				Level:        "",
				OutputWriter: "",
			},
			Expected: Expected{
				Level:            LevelOff,
				OutputWriterName: CLILogOutputWriteName,
			},
		},
		{
			Data: config.Logger{
				Level:        string(LevelOff),
				OutputWriter: string(cli),
			},
			Expected: Expected{
				Level:            LevelOff,
				OutputWriterName: CLILogOutputWriteName,
			},
		},
		{
			Data: config.Logger{
				Level:        string(LevelInfo),
				OutputWriter: string(jsonFile),
			},
			Expected: Expected{
				Level:            LevelInfo,
				OutputWriterName: JSONLogOutputFileWriteName,
			},
		},
		{
			Data: config.Logger{
				Level:        string(LevelWarning),
				OutputWriter: string(textFile),
			},
			Expected: Expected{
				Level:            LevelWarning,
				OutputWriterName: TextLogOutputFileWriteName,
			},
		},
		{
			Data: config.Logger{
				Level:        string(LevelError),
				OutputWriter: "undefined",
			},
			Expected: Expected{
				Level:            LevelError,
				OutputWriterName: CLILogOutputWriteName,
			},
		},
		{
			Data: config.Logger{
				Level:        string(LevelDebug),
				OutputWriter: "undefined",
			},
			Expected: Expected{
				Level:            LevelDebug,
				OutputWriterName: CLILogOutputWriteName,
			},
		},
	}

	// config.Init()
	config.LoggerConfig = &config.Config{}
	for idx, test := range testCases {
		loggerOnce = sync.Once{}
		config.LoggerConfig.Logger = test.Data
		log := NewLog()
		assert.Equal(t, test.Expected.Level, log.loggerLevel)
		fnName, err := itesting.GetFunctionName(log.outputWrite)
		if err != nil {
			t.FailNow()
		}
		assert.Contains(t, fnName, test.Expected.OutputWriterName, fmt.Sprintf("Test #%d: [level=%s] [outputWriter=%s]\n", idx+1, test.Data.Level, test.Data.OutputWriter))
	}
}

func TestWithLoggerLevel(t *testing.T) {
	l := logging{}
	WithLoggerLevel(LevelDebug)(&l)
	assert.Equal(t, LevelDebug, l.loggerLevel)
}

func TestWithLogOutputWriter(t *testing.T) {
	CLI := CLILogOutputWrite()
	l := logging{}
	WithLogOutputWriter(CLI)(&l)
	fnName, err := itesting.GetFunctionName(l.outputWrite)
	if err != nil {
		t.FailNow()
	}
	assert.Contains(t, fnName, "CLILogOutputWrite")
}

func TestLoggingOff(t *testing.T) {
	loggerOnce = sync.Once{}
	log := NewLog(WithLoggerLevel(LevelOff))

	bytes, err := itesting.CaptureOutput(func() error {
		log.Info("test info", nil)
		log.Warning("test warning", nil)
		log.Error("test error", nil)
		log.Debug("test debug", nil)
		return nil
	})
	if err != nil {
		t.Errorf("error: could not capture stdout output %v", err)
		return
	}

	assert.NotContains(t, bytes, "test")
}

func TestInfo(t *testing.T) {
	loggerOnce = sync.Once{}
	log := NewLog(WithLoggerLevel(LevelInfo))

	bytes, err := itesting.CaptureOutput(func() error {
		log.Info("test info", map[string]any{
			"varInt":   0,
			"varStr":   "string",
			"varFloat": 3.14,
		})
		log.Warning("test warning", nil)
		log.Error("test error", nil)
		log.Debug("test debug", nil)
		return nil
	})
	if err != nil {
		t.Errorf("error: could not capture stdout output %v", err)
		return
	}

	assert.Contains(t, bytes, "varInt=0")
	assert.Contains(t, bytes, "varStr=string")
	assert.Contains(t, bytes, "varFloat=3.14")
	assert.Contains(t, bytes, "test info")
	assert.NotContains(t, bytes, "test warning")
	assert.NotContains(t, bytes, "test error")
	assert.NotContains(t, bytes, "test debug")
}

func TestWarning(t *testing.T) {
	loggerOnce = sync.Once{}
	log := NewLog(WithLoggerLevel(LevelWarning))

	bytes, err := itesting.CaptureOutput(func() error {
		log.Warning("test warning", map[string]any{
			"varInt":   0,
			"varStr":   "string",
			"varFloat": 3.14,
		})
		log.Info("test info", nil)
		log.Error("test error", nil)
		log.Debug("test debug", nil)
		return nil
	})
	if err != nil {
		t.Errorf("error: could not capture stdout output %v", err)
		return
	}

	assert.Contains(t, bytes, "varInt=0")
	assert.Contains(t, bytes, "varStr=string")
	assert.Contains(t, bytes, "varFloat=3.14")
	assert.Contains(t, bytes, "test info")
	assert.Contains(t, bytes, "test warning")
	assert.NotContains(t, bytes, "test error")
	assert.NotContains(t, bytes, "test debug")
}

func TestError(t *testing.T) {
	loggerOnce = sync.Once{}
	log := NewLog(WithLoggerLevel(LevelError))

	bytes, err := itesting.CaptureOutput(func() error {
		log.Error("test error", map[string]any{
			"varInt":   0,
			"varStr":   "string",
			"varFloat": 3.14,
		})
		log.Info("test info", nil)
		log.Warning("test warning", nil)
		log.Debug("test debug", nil)
		return nil
	})
	if err != nil {
		t.Errorf("error: could not capture stdout output %v", err)
		return
	}

	assert.Contains(t, bytes, "varInt=0")
	assert.Contains(t, bytes, "varStr=string")
	assert.Contains(t, bytes, "varFloat=3.14")
	assert.Contains(t, bytes, "test info")
	assert.Contains(t, bytes, "test warning")
	assert.Contains(t, bytes, "test error")
	assert.NotContains(t, bytes, "test debug")
}

func TestDebug(t *testing.T) {
	loggerOnce = sync.Once{}
	log := NewLog(WithLoggerLevel(LevelDebug))

	bytes, err := itesting.CaptureOutput(func() error {
		log.Debug("test debug", map[string]any{
			"varInt":   0,
			"varStr":   "string",
			"varFloat": 3.14,
		})
		log.Info("test info", nil)
		log.Warning("test warning", nil)
		log.Error("test error", nil)
		return nil
	})
	if err != nil {
		t.Errorf("error: could not capture stdout output %v", err)
		return
	}

	assert.Contains(t, bytes, "varInt=0")
	assert.Contains(t, bytes, "varStr=string")
	assert.Contains(t, bytes, "varFloat=3.14")
	assert.Contains(t, bytes, "test info")
	assert.Contains(t, bytes, "test warning")
	assert.Contains(t, bytes, "test error")
	assert.Contains(t, bytes, "test debug")
}
