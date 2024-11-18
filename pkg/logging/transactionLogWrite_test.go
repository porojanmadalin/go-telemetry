package logging

import (
	"encoding/json"
	"fmt"
	"go-telemetry/pkg/internal/config"
	itesting "go-telemetry/pkg/internal/telemetrytesting"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var testTransactionLogging = TransactionLoggerData{
	LoggerLevel: LevelError,
	TransactionLogs: []*LoggerData{
		{
			LoggerLevel: LevelInfo,
			Timestamp:   now,
			Message:     "test1",
			MetaData: map[string]any{
				"varInt":   1,
				"varStr":   "string1",
				"varFloat": 3.14,
			},
		},
		{
			LoggerLevel: LevelWarning,
			Timestamp:   now.Add(1 * time.Second),
			Message:     "test2",
			MetaData: map[string]any{
				"varInt":   2,
				"varStr":   "string2",
				"varFloat": 6.28,
			},
		},
	},
}

func transactionLogFormatStart(transactionId string, startTimestamp time.Time) string {
	return fmt.Sprintf("[%s] Transaction {%s} started!", startTimestamp.Format(timestampFormat), transactionId)
}
func transactionLogFormatEnd(transactionId string, endTimestamp time.Time) string {
	return fmt.Sprintf("[%s] Transaction {%s} ended!", endTimestamp.Format(timestampFormat), transactionId)
}

func TestCLITransactionLogOutputWrite(t *testing.T) {
	setupTestEnvironment(t, transactionLogTestDirName)

	output, err := itesting.CaptureOutput(func() error {
		err := CLITransactionLogOutputWrite()(testTransactionId, now, testEndTimestamp, &testTransactionLogging)
		return err
	})
	if err != nil {
		t.Fatalf("fatal: could not capture stdout output %v", err)
	}

	outputStr := string(output[:])

	outputLines := strings.Split(outputStr, "\n")

	assert.True(t, len(testTransactionLogging.TransactionLogs) == len(outputLines)-3)
	assert.Equal(t, transactionLogFormatStart(testTransactionId, now), outputLines[0])
	assert.Equal(t, transactionLogFormatEnd(testTransactionId, testEndTimestamp), outputLines[len(outputLines)-2])

	for i := 1; i <= len(outputLines)-3; i++ {
		assert.Contains(t, outputLines[i], "--> "+logFormat(*testTransactionLogging.TransactionLogs[i-1]))
		assert.Contains(t, outputLines[i], fmt.Sprintf("varInt=%d", testTransactionLogging.TransactionLogs[i-1].MetaData["varInt"]))
		assert.Contains(t, outputLines[i], fmt.Sprintf("varStr=%s", testTransactionLogging.TransactionLogs[i-1].MetaData["varStr"]))
		assert.Contains(t, outputLines[i], fmt.Sprintf("varFloat=%f", testTransactionLogging.TransactionLogs[i-1].MetaData["varFloat"]))
	}
}

func TestJSONTransactionLogOutputFileWrite(t *testing.T) {
	setupTestEnvironment(t, transactionLogTestDirName)

	generatedFileName := fmt.Sprintf("%s.json", now.Format(fileTimestampFormat))
	t.Cleanup(func() {
		cleanup(t, []string{generatedFileName})
	})

	err := JSONTransactionLogOutputFileWrite()(testTransactionId, now, testEndTimestamp, &testTransactionLogging)
	if err != nil {
		t.Fatalf("fatal: could not write logs to text file %v", err)
	}

	f, err := os.OpenFile(filepath.Join(config.LoggerConfig.Logger.OutputDir, generatedFileName), os.O_RDONLY, 0644)
	if err != nil {
		t.Fatalf("error: could not open json file %v", err)
	}
	defer f.Close()

	var logs []struct {
		TransactionId         string    `json:"transactionId"`
		StartTimestamp        time.Time `json:"startTimestamp"`
		EndTimestamp          time.Time `json:"endTimestamp"`
		TransactionLoggerData `json:"transactionData"`
	}

	jsonParser := json.NewDecoder(f)
	if err = jsonParser.Decode(&logs); err != nil {
		t.Fatalf("error: could not decode json file %v", err)
	}

	assert.Equal(t, testTransactionLogging.LoggerLevel, logs[0].LoggerLevel)

	for i := range len(testTransactionLogging.TransactionLogs) {
		assert.True(t, testTransactionLogging.TransactionLogs[i].Timestamp.Equal(logs[0].TransactionLogs[i].Timestamp))
		assert.Equal(t, testTransactionLogging.TransactionLogs[i].Message, logs[0].TransactionLogs[i].Message)
		assert.Equal(t, testTransactionLogging.TransactionLogs[i].MetaData["varInt"], int(logs[0].TransactionLogs[i].MetaData["varInt"].(float64))) // JSON numbers are parsed as floats
		assert.Equal(t, testTransactionLogging.TransactionLogs[i].MetaData["varStr"], logs[0].TransactionLogs[i].MetaData["varStr"])
		assert.True(t, itesting.AlmostEqual(logs[0].TransactionLogs[i].MetaData["varFloat"].(float64), testTransactionLogging.TransactionLogs[i].MetaData["varFloat"].(float64)))
	}
}

func TestTextTransactionLogOutputFileWrite(t *testing.T) {
	setupTestEnvironment(t, transactionLogTestDirName)

	now := time.Now()
	generatedFileName := fmt.Sprintf("%s.log", now.Format(fileTimestampFormat))
	t.Cleanup(func() {
		cleanup(t, []string{generatedFileName})
	})

	err := TextTransactionLogOutputFileWrite()(testTransactionId, now, testEndTimestamp, &testTransactionLogging)
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

	outputStr := string(b[:])

	outputLines := strings.Split(outputStr, "\n")

	assert.True(t, len(testTransactionLogging.TransactionLogs) == len(outputLines)-3)
	assert.Equal(t, transactionLogFormatStart(testTransactionId, now), outputLines[0])
	assert.Equal(t, transactionLogFormatEnd(testTransactionId, testEndTimestamp), outputLines[len(outputLines)-2])

	for i := 1; i <= len(outputLines)-3; i++ {
		assert.Contains(t, outputLines[i], "--> "+logFormat(*testTransactionLogging.TransactionLogs[i-1]))
		assert.Contains(t, outputLines[i], fmt.Sprintf("varInt=%d", testTransactionLogging.TransactionLogs[i-1].MetaData["varInt"]))
		assert.Contains(t, outputLines[i], fmt.Sprintf("varStr=%s", testTransactionLogging.TransactionLogs[i-1].MetaData["varStr"]))
		assert.Contains(t, outputLines[i], fmt.Sprintf("varFloat=%f", testTransactionLogging.TransactionLogs[i-1].MetaData["varFloat"]))
	}
}
