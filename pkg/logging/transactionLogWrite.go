package logging

import (
	"encoding/json"
	"fmt"
	"go-telemetry/pkg/internal/config"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"
)

// A TransactionLogOutputWriter is a output writer function for transaction logging.
type TransactionLogOutputWriter func(transactionId string, startTimestamp time.Time, endTimestamp time.Time, transactionLoggerData *TransactionLoggerData) error

// CLITransactionLogOutputWrite returns an output writer that prints the transaction log to the CLI.
func CLITransactionLogOutputWrite() TransactionLogOutputWriter {
	return func(transactionId string, startTimestamp time.Time, endTimestamp time.Time, transactionLoggerData *TransactionLoggerData) error {
		fmt.Printf("[%s] Transaction {%s} started!\n", startTimestamp.Format(timestampFormat), transactionId)

		for _, entry := range transactionLoggerData.TransactionLogs {
			fmt.Printf("--> [%s] [%s] %s", entry.Timestamp.Format(timestampFormat), entry.LoggerLevel, entry.Message)
			for k, v := range entry.MetaData {
				typeName := reflect.TypeOf(v).Name()
				if strings.Contains(typeName, "int") {
					fmt.Printf(" [%s=%d]", k, v)
				} else if strings.Contains(typeName, "float") {
					fmt.Printf(" [%s=%f]", k, v)
				} else if typeName == "string" {
					fmt.Printf(" [%s=%s]", k, v)
				} else {
					fmt.Printf("\n")
					return fmt.Errorf("error: formatting for this variable type is not implement %s", typeName)
				}
			}
			fmt.Printf("\n")
		}

		fmt.Printf("[%s] Transaction {%s} ended!\n", endTimestamp.Format(timestampFormat), transactionId)
		return nil
	}
}

// JSONTransactionLogOutputFileWrite returns an output writer that prints the transaction log to a JSON file.
func JSONTransactionLogOutputFileWrite() TransactionLogOutputWriter {
	return func(transactionId string, startTimestamp time.Time, endTimestamp time.Time, transactionLoggerData *TransactionLoggerData) error {
		fileName := fmt.Sprintf(filepath.Join(config.LoggerConfig.Logger.OutputDir, "%s_transactions.json"), time.Now().Format(fileTimestampFormat))
		f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			return fmt.Errorf("error: could not open json file %v", err)
		}
		defer f.Close()

		ret, err := f.Seek(-int64(len([]byte(endArray))), io.SeekEnd)
		if err != nil {
			_, err := f.Write([]byte(startArray))
			if err != nil {
				return fmt.Errorf("error: could not write json file %v", err)
			}
			ret = int64(len([]byte(startArray)))
		}

		type OutputJSON struct {
			TransactionID   string                 `json:"transactionId"`
			StartTimestamp  time.Time              `json:"startTimestamp"`
			EndTimestamp    time.Time              `json:"endTimestamp"`
			TransactionLogs *TransactionLoggerData `json:"transactionData"`
		}
		outputJSON := &OutputJSON{
			TransactionID:   transactionId,
			StartTimestamp:  startTimestamp,
			EndTimestamp:    endTimestamp,
			TransactionLogs: transactionLoggerData,
		}

		loggerDataBytes, err := json.MarshalIndent(outputJSON, indent, indent)
		if err != nil {
			return fmt.Errorf("error: could not marshal logger data %v", err)
		}

		writtenDataBytes := []byte(objectDelimiter)
		if ret == int64(len([]byte(startArray))) {
			writtenDataBytes = []byte(indent)
		}
		writtenDataBytes = append(writtenDataBytes, loggerDataBytes...)

		_, err = f.WriteAt(writtenDataBytes, ret)
		if err != nil {
			return fmt.Errorf("error: could not write json file %v", err)
		}

		_, err = f.WriteAt([]byte(endArray), ret+int64(len(writtenDataBytes)))
		if err != nil {
			return fmt.Errorf("error: could not write json file %v", err)
		}

		return nil
	}
}

// TextTransactionLogOutputFileWrite returns an output writer that prints the transaction log to a text file.
func TextTransactionLogOutputFileWrite() TransactionLogOutputWriter {
	return func(transactionId string, startTimestamp time.Time, endTimestamp time.Time, transactionLoggerData *TransactionLoggerData) error {
		fileName := fmt.Sprintf(filepath.Join(config.LoggerConfig.Logger.OutputDir, "%s_transactions.json"), time.Now().Format(fileTimestampFormat))
		f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			return fmt.Errorf("error: could not open text file %v", err)
		}
		defer f.Close()

		f.WriteString(fmt.Sprintf("[%s] Transaction {%s} started!\n", startTimestamp.Format(timestampFormat), transactionId))

		for _, entry := range transactionLoggerData.TransactionLogs {
			f.WriteString(fmt.Sprintf("--> [%s] [%s] %s", entry.Timestamp.Format(timestampFormat), entry.LoggerLevel, entry.Message))
			for k, v := range entry.MetaData {
				typeName := reflect.TypeOf(v).Name()
				if strings.Contains(typeName, "int") {
					f.WriteString(fmt.Sprintf(" [%s=%d]", k, v))
				} else if strings.Contains(typeName, "float") {
					f.WriteString(fmt.Sprintf(" [%s=%f]", k, v))
				} else if typeName == "string" {
					f.WriteString(fmt.Sprintf(" [%s=%s]", k, v))
				} else {
					f.WriteString("\n")
					return fmt.Errorf("error: formatting for this variable type is not implement %s", typeName)
				}
			}
			f.WriteString("\n")
		}

		f.WriteString(fmt.Sprintf("[%s] Transaction {%s} ended!\n", endTimestamp.Format(timestampFormat), transactionId))
		return nil
	}
}
