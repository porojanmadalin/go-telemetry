package logging

import (
	"encoding/json"
	"fmt"
	"go-telemetry/pkg/internal/config"
	itesting "go-telemetry/pkg/internal/telemetrytesting"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const logTestDirName = "testlog"

var now = time.Now()
var testLogging = LoggerData{
	LoggerLevel: LevelInfo,
	Timestamp:   now,
	Message:     "test",
	MetaData: map[string]any{
		"varInt":   0,
		"varStr":   "string",
		"varFloat": 3.14,
	},
}

func logFormat(loggerData LoggerData) string {
	return fmt.Sprintf("[%s] [%s] %s", loggerData.Timestamp.Format(timestampFormat), loggerData.LoggerLevel, loggerData.Message)
}

func TestCLILogOutputWrite(t *testing.T) {
	setupTestEnvironment(t, logTestDirName)

	output, err := itesting.CaptureOutput(func() error {
		err := CLILogOutputWrite()(&testLogging)
		return err
	})
	if err != nil {
		t.Fatalf("fatal: could not capture stdout output %v", err)
	}

	assert.Contains(t, output, logFormat(testLogging))
	assert.Contains(t, output, fmt.Sprintf("varInt=%d", testLogging.MetaData["varInt"]))
	assert.Contains(t, output, fmt.Sprintf("varStr=%s", testLogging.MetaData["varStr"]))
	assert.Contains(t, output, fmt.Sprintf("varFloat=%f", testLogging.MetaData["varFloat"]))
	assert.Contains(t, output, "\n")
}

func TestJSONLogOutputFileWrite(t *testing.T) {
	setupTestEnvironment(t, logTestDirName)

	generatedFileName := fmt.Sprintf("%s.json", now.Format(fileTimestampFormat))
	t.Cleanup(func() {
		cleanup(t, []string{generatedFileName})
	})

	err := JSONLogOutputFileWrite()(&testLogging)
	if err != nil {
		t.Fatalf("fatal: could not write logs to text file %v", err)
	}

	f, err := os.OpenFile(filepath.Join(config.LoggerConfig.Logger.OutputDir, generatedFileName), os.O_RDONLY, 0644)
	if err != nil {
		t.Fatalf("error: could not open json file %v", err)
	}
	defer f.Close()

	var logs []LoggerData

	jsonParser := json.NewDecoder(f)
	if err = jsonParser.Decode(&logs); err != nil {
		t.Fatalf("error: could not decode json file %v", err)
	}

	assert.True(t, testLogging.Timestamp.Equal(logs[0].Timestamp))
	assert.Equal(t, testLogging.LoggerLevel, logs[0].LoggerLevel)
	assert.Equal(t, testLogging.Message, logs[0].Message)
	assert.Equal(t, testLogging.MetaData["varInt"], int(logs[0].MetaData["varInt"].(float64))) // JSON numbers are parsed as floats
	assert.Equal(t, testLogging.MetaData["varStr"], logs[0].MetaData["varStr"])
	assert.True(t, itesting.AlmostEqual(logs[0].MetaData["varFloat"].(float64), testLogging.MetaData["varFloat"].(float64)))
}

func TestTextLogOutputFileWrite(t *testing.T) {
	setupTestEnvironment(t, logTestDirName)

	generatedFileName := fmt.Sprintf("%s.log", now.Format(fileTimestampFormat))
	t.Cleanup(func() {
		cleanup(t, []string{generatedFileName})
	})

	err := TextLogOutputFileWrite()(&testLogging)
	if err != nil {
		t.Fatalf("fatal: could not write logs to text file %v", err)
	}

	f, err := os.OpenFile(filepath.Join(config.LoggerConfig.Logger.OutputDir, generatedFileName), os.O_RDONLY, 0644)
	if err != nil {
		t.Fatalf("error: could not open text file %v", err)
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		t.Fatalf("error: could not read from text file %v", err)
	}

	output := string(b[:])

	assert.Contains(t, output, logFormat(testLogging))
	assert.Contains(t, output, fmt.Sprintf("varInt=%d", testLogging.MetaData["varInt"]))
	assert.Contains(t, output, fmt.Sprintf("varStr=%s", testLogging.MetaData["varStr"]))
	assert.Contains(t, output, fmt.Sprintf("varFloat=%f", testLogging.MetaData["varFloat"]))
	assert.Contains(t, output, "\n")
}
