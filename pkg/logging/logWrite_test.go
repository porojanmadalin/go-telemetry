package logging

import (
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

const (
	testOutputDir = ""
)

func logFormat(loggerData LoggerData) string {
	return fmt.Sprintf("[%s] [%s] %s", loggerData.Timestamp.Format(timestampFormat), loggerData.LoggerLevel, loggerData.Message)
}

func TestCLILogOutputWrite(t *testing.T) {
	setupTestEnvironment(t)

	now := time.Now()
	testLogging := LoggerData{
		LoggerLevel: LevelInfo,
		Timestamp:   now,
		Message:     "test",
		MetaData: map[string]any{
			"varInt":   0,
			"varStr":   "string",
			"varFloat": 3.14,
		},
	}

	output, err := itesting.CaptureOutput(func() error {
		err := CLILogOutputWrite()(&testLogging)
		return err
	})
	if err != nil {
		t.Fatalf("fatal: could not capture stdout output %v", err)
	}

	assert.Contains(t, output, logFormat(testLogging))
	assert.Contains(t, output, "varInt=0")
	assert.Contains(t, output, "varStr=string")
	assert.Contains(t, output, "varFloat=3.14")
	assert.Contains(t, output, "\n")
}

func TestJSONLogOutputFileWrite(t *testing.T) {
	setupTestEnvironment(t)

}

func TestTextLogOutputFileWrite(t *testing.T) {
	setupTestEnvironment(t)

	now := time.Now()
	generatedFileName := fmt.Sprintf("%s.log", now.Format(fileTimestampFormat))
	t.Cleanup(func() {
		cleanup(t, []string{generatedFileName})
	})

	testLogging := LoggerData{
		LoggerLevel: LevelInfo,
		Timestamp:   now,
		Message:     "test",
		MetaData: map[string]any{
			"varInt":   0,
			"varStr":   "string",
			"varFloat": 3.14,
		},
	}

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
		t.Fatalf("error: could not open text file %v", err)
	}

	output := string(b[:])

	assert.Contains(t, output, logFormat(testLogging))
	assert.Contains(t, output, "varInt=0")
	assert.Contains(t, output, "varStr=string")
	assert.Contains(t, output, "varFloat=3.14")
	assert.Contains(t, output, "\n")
}
