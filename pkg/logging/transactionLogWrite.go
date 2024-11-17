package logging

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"time"
)

type TransactionLogOutputWriter func(transactionId string, startTimestamp time.Time, endTimestamp time.Time, transactionLoggerData *TransactionLoggerData) error

func CLITransactionLogOutputWrite() TransactionLogOutputWriter {
	return func(transactionId string, startTimestamp time.Time, endTimestamp time.Time, transactionLoggerData *TransactionLoggerData) error {
		fmt.Printf("[%s] Transaction {%s} started!\n", startTimestamp.Format("2006-01-02 15:04:05.0000"), transactionId)

		for _, entry := range transactionLoggerData.TransactionLogs {
			fmt.Printf("--> [%s] [%s] %s", entry.Timestamp.Format("2006-01-02 15:04:05.0000"), entry.LoggerLevel, entry.Message)
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

		fmt.Printf("[%s] Transaction {%s} ended!\n", endTimestamp.Format("2006-01-02 15:04:05.0000"), transactionId)
		return nil
	}
}

func JSONTransactionLogOutputFileWrite() TransactionLogOutputWriter {
	return func(transactionId string, startTimestamp time.Time, endTimestamp time.Time, transactionLoggerData *TransactionLoggerData) error {
		fileName := fmt.Sprintf("%s.json", time.Now().Format("2006-01-02"))
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
			TransactionLogs *TransactionLoggerData `json:"transactionLogs"`
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

func TextTransactionLogOutputFileWrite() TransactionLogOutputWriter {
	return func(transactionId string, startTimestamp time.Time, endTimestamp time.Time, transactionLoggerData *TransactionLoggerData) error {
		fileName := fmt.Sprintf("%s.log", time.Now().Format("2006-01-02"))
		f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			return fmt.Errorf("error: could not open text file %v", err)
		}
		defer f.Close()

		f.WriteString(fmt.Sprintf("[%s] Transaction {%s} started!\n", startTimestamp.Format("2006-01-02 15:04:05.0000"), transactionId))

		for _, entry := range transactionLoggerData.TransactionLogs {
			f.WriteString(fmt.Sprintf("--> [%s] [%s] %s", entry.Timestamp.Format("2006-01-02 15:04:05.0000"), entry.LoggerLevel, entry.Message))
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

		f.WriteString(fmt.Sprintf("[%s] Transaction {%s} ended!\n", endTimestamp.Format("2006-01-02 15:04:05.0000"), transactionId))
		return nil
	}
}
