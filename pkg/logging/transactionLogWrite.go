package logging

import (
	"fmt"
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
		// ! Not implemented
		return nil
	}
}

func TextTransactionLogOutputFileWrite() TransactionLogOutputWriter {
	return func(transactionId string, startTimestamp time.Time, endTimestamp time.Time, transactionLoggerData *TransactionLoggerData) error {
		fileName := fmt.Sprintf("%s.log", startTimestamp.Format("2006-01-02"))
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
