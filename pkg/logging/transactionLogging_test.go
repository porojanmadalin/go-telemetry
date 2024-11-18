package logging

import (
	"go-telemetry/pkg/internal/config"
	itesting "go-telemetry/pkg/internal/telemetrytesting"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	// output writers fn names
	CLITransactionLogOutputWriteName      = "CLITransactionLogOutputWrite"
	JSONTransactionLogOutputFileWriteName = "JSONTransactionLogOutputFileWrite"
	TextTransactionLogOutputFileWriteName = "TextTransactionLogOutputFileWrite"
)

const transactionLogTestDirName = "testtransactionlog"

var testEndTimestamp = now.Add(1 * time.Second)
var testTransactionId = "testTransaction"

func TestNewTransactionLogWithYAMLConfigTableDriven(t *testing.T) {
	type Expected struct {
		Level            loggerLevel
		OutputWriterName string
	}
	type TestCase struct {
		TestName string
		Data     config.Logger
		Expected Expected
	}

	testCases := []TestCase{
		{
			TestName: "Config contains invalid values",
			Data: config.Logger{
				Level:        "undefined",
				OutputWriter: "undefined",
			},
			Expected: Expected{
				Level:            LevelInfo,
				OutputWriterName: CLITransactionLogOutputWriteName,
			},
		},
		{
			TestName: "Unset Config",
			Data: config.Logger{
				Level:        "",
				OutputWriter: "",
			},
			Expected: Expected{
				Level:            LevelInfo,
				OutputWriterName: CLITransactionLogOutputWriteName,
			},
		},
		{
			TestName: "Log level off, output writer cli",
			Data: config.Logger{
				Level:        string(LevelOff),
				OutputWriter: string(cli),
			},
			Expected: Expected{
				Level:            LevelOff,
				OutputWriterName: CLITransactionLogOutputWriteName,
			},
		},
		{
			TestName: "Log level info, output writer json file",
			Data: config.Logger{
				Level:        string(LevelInfo),
				OutputWriter: string(jsonFile),
			},
			Expected: Expected{
				Level:            LevelInfo,
				OutputWriterName: JSONTransactionLogOutputFileWriteName,
			},
		},
		{
			TestName: "Log level warning, output writer text file",
			Data: config.Logger{
				Level:        string(LevelWarning),
				OutputWriter: string(textFile),
			},
			Expected: Expected{
				Level:            LevelWarning,
				OutputWriterName: TextTransactionLogOutputFileWriteName,
			},
		},
		{
			TestName: "Log level error, output writer invalid value",
			Data: config.Logger{
				Level:        string(LevelError),
				OutputWriter: "undefined",
			},
			Expected: Expected{
				Level:            LevelError,
				OutputWriterName: CLITransactionLogOutputWriteName,
			},
		},
		{
			TestName: "Log level debug, output writer invalid value",
			Data: config.Logger{
				Level:        string(LevelDebug),
				OutputWriter: "undefined",
			},
			Expected: Expected{
				Level:            LevelDebug,
				OutputWriterName: CLITransactionLogOutputWriteName,
			},
		},
	}

	config.LoggerConfig = &config.Config{}
	for _, test := range testCases {
		t.Run(test.TestName, func(t *testing.T) {
			transactionLoggerOnce = sync.Once{}
			config.LoggerConfig.Logger = test.Data
			log, err := NewTransactionLog(testTransactionId)
			if err != nil {
				t.Fatalf("fatal: transaction log could not be initialized for transaction %s %v", testTransactionId, err)
			}
			assert.Equal(t, test.Expected.Level, log.loggerLevel)
			fnName, err := itesting.GetFunctionName(log.outputWrite)
			if err != nil {
				t.Fatalf("fatal: could not locate the function %v", err)
			}
			assert.Contains(t, fnName, test.Expected.OutputWriterName)
		})
	}
}

func TestWithTransactionLoggerLevel(t *testing.T) {
	l := transactionLogging{}
	WithTransactionLoggerLevel(LevelDebug)(&l)
	assert.Equal(t, LevelDebug, l.loggerLevel)
}

func TestWithTransactionLogOutputWriter(t *testing.T) {
	CLI := CLITransactionLogOutputWrite()
	l := transactionLogging{}
	WithTransactionLogOutputWriter(CLI)(&l)
	fnName, err := itesting.GetFunctionName(l.outputWrite)
	if err != nil {
		t.Fatalf("fatal: could not locate the function %v", err)
	}
	assert.Contains(t, fnName, CLITransactionLogOutputWriteName)
}

func TestTransactionLoggingOff(t *testing.T) {
	transactionLoggerOnce = sync.Once{}
	log, err := NewTransactionLog(testTransactionId, WithTransactionLoggerLevel(LevelOff))
	if err != nil {
		t.Fatalf("fatal: transaction log could not be initialized for transaction %s %v", testTransactionId, err)
	}

	bytes, err := itesting.CaptureOutput(func() error {
		err := log.StartTransactionLogging()
		if err != nil {
			return err
		}
		log.Info("test info", nil)
		log.Warning("test warning", nil)
		log.Error("test error", nil)
		log.Debug("test debug", nil)
		err = log.StopTransactionLogging()
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		t.Errorf("error: could not capture stdout output %v", err)
		return
	}

	assert.NotContains(t, bytes, "test")
	assert.NotContains(t, bytes, testTransactionId)
	assert.NotContains(t, bytes, "started")
	assert.NotContains(t, bytes, "ended")
}

func TestTransactionInfo(t *testing.T) {
	transactionLoggerOnce = sync.Once{}
	log, err := NewTransactionLog(testTransactionId, WithTransactionLoggerLevel(LevelInfo))
	if err != nil {
		t.Fatalf("fatal: transaction log could not be initialized for transaction %s %v", testTransactionId, err)
	}

	bytes, err := itesting.CaptureOutput(func() error {
		err := log.StartTransactionLogging()
		if err != nil {
			return err
		}
		log.Info("test info", map[string]any{
			"varInt":   0,
			"varStr":   "string",
			"varFloat": 3.14,
		})
		log.Warning("test warning", nil)
		log.Error("test error", nil)
		log.Debug("test debug", nil)
		err = log.StopTransactionLogging()
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		t.Fatalf("error: could not capture stdout output %v", err)
	}

	assert.Contains(t, bytes, testTransactionId)
	assert.Contains(t, bytes, "started")
	assert.Contains(t, bytes, "ended")
	assert.Contains(t, bytes, "varInt=0")
	assert.Contains(t, bytes, "varStr=string")
	assert.Contains(t, bytes, "varFloat=3.14")
	assert.Contains(t, bytes, "test info")
	assert.NotContains(t, bytes, "test warning")
	assert.NotContains(t, bytes, "test error")
	assert.NotContains(t, bytes, "test debug")
}

func TestTransactionWarning(t *testing.T) {
	transactionLoggerOnce = sync.Once{}
	log, err := NewTransactionLog(testTransactionId, WithTransactionLoggerLevel(LevelWarning))
	if err != nil {
		t.Fatalf("fatal: transaction log could not be initialized for transaction %s %v", testTransactionId, err)
	}

	bytes, err := itesting.CaptureOutput(func() error {
		err := log.StartTransactionLogging()
		if err != nil {
			return err
		}
		log.Warning("test warning", map[string]any{
			"varInt":   0,
			"varStr":   "string",
			"varFloat": 3.14,
		})
		log.Info("test info", nil)
		log.Error("test error", nil)
		log.Debug("test debug", nil)
		err = log.StopTransactionLogging()
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		t.Fatalf("error: could not capture stdout output %v", err)
	}

	assert.Contains(t, bytes, testTransactionId)
	assert.Contains(t, bytes, "started")
	assert.Contains(t, bytes, "ended")
	assert.Contains(t, bytes, "varInt=0")
	assert.Contains(t, bytes, "varStr=string")
	assert.Contains(t, bytes, "varFloat=3.14")
	assert.Contains(t, bytes, "test info")
	assert.Contains(t, bytes, "test warning")
	assert.NotContains(t, bytes, "test error")
	assert.NotContains(t, bytes, "test debug")
}

func TestTransactionError(t *testing.T) {
	transactionLoggerOnce = sync.Once{}
	log, err := NewTransactionLog(testTransactionId, WithTransactionLoggerLevel(LevelError))
	if err != nil {
		t.Fatalf("fatal: transaction log could not be initialized for transaction %s %v", testTransactionId, err)
	}

	bytes, err := itesting.CaptureOutput(func() error {
		err := log.StartTransactionLogging()
		if err != nil {
			return err
		}
		log.Error("test error", map[string]any{
			"varInt":   0,
			"varStr":   "string",
			"varFloat": 3.14,
		})
		log.Info("test info", nil)
		log.Warning("test warning", nil)
		log.Debug("test debug", nil)
		err = log.StopTransactionLogging()
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		t.Fatalf("error: could not capture stdout output %v", err)
	}

	assert.Contains(t, bytes, testTransactionId)
	assert.Contains(t, bytes, "started")
	assert.Contains(t, bytes, "ended")
	assert.Contains(t, bytes, "varInt=0")
	assert.Contains(t, bytes, "varStr=string")
	assert.Contains(t, bytes, "varFloat=3.14")
	assert.Contains(t, bytes, "test info")
	assert.Contains(t, bytes, "test warning")
	assert.Contains(t, bytes, "test error")
	assert.NotContains(t, bytes, "test debug")
}

func TestTransactionDebug(t *testing.T) {
	transactionLoggerOnce = sync.Once{}
	log, err := NewTransactionLog(testTransactionId, WithTransactionLoggerLevel(LevelDebug))
	if err != nil {
		t.Fatalf("fatal: transaction log could not be initialized for transaction %s %v", testTransactionId, err)
	}

	bytes, err := itesting.CaptureOutput(func() error {
		err := log.StartTransactionLogging()
		if err != nil {
			return err
		}
		log.Debug("test debug", map[string]any{
			"varInt":   0,
			"varStr":   "string",
			"varFloat": 3.14,
		})
		log.Info("test info", nil)
		log.Warning("test warning", nil)
		log.Error("test error", nil)
		err = log.StopTransactionLogging()
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		t.Fatalf("error: could not capture stdout output %v", err)
	}

	assert.Contains(t, bytes, testTransactionId)
	assert.Contains(t, bytes, "started")
	assert.Contains(t, bytes, "ended")
	assert.Contains(t, bytes, "varInt=0")
	assert.Contains(t, bytes, "varStr=string")
	assert.Contains(t, bytes, "varFloat=3.14")
	assert.Contains(t, bytes, "test info")
	assert.Contains(t, bytes, "test warning")
	assert.Contains(t, bytes, "test error")
	assert.Contains(t, bytes, "test debug")
}

func TestStartTransactionLogging(t *testing.T) {

}

func StopTransactionLogging(t *testing.T) {

}
