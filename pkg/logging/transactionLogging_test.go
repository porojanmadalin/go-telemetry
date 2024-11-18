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

func TestNewTransactionLogWithYAMLConfig(t *testing.T) {
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
				Level:            LevelInfo,
				OutputWriterName: CLITransactionLogOutputWriteName,
			},
		},
		{
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
		transactionLoggerOnce = sync.Once{}
		config.LoggerConfig.Logger = test.Data
		log, err := NewTransactionLog(testTransactionId)
		if err != nil {
			t.Fatalf("fatal: transaction log could not be initialized for transaction %s %v", testTransactionId, err)
		}
		assert.Equal(t, test.Expected.Level, log.loggerLevel)
		fnName, err := itesting.GetFunctionName(log.outputWrite)
		if err != nil {
			t.FailNow()
		}
		assert.Contains(t, fnName, test.Expected.OutputWriterName)
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
		t.FailNow()
	}
	assert.Contains(t, fnName, CLITransactionLogOutputWriteName)
}
