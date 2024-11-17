package logging

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

type TransactionLogOutputWriter func(transactionId string, startTimestamp time.Time, endTimestamp time.Time, transactionLoggerData *TransactionLoggerData) error

func CLITransactionLogOutputWrite() TransactionLogOutputWriter {
	return func(transactionId string, startTimestamp time.Time, endTimeStamp time.Time, transactionLoggerData *TransactionLoggerData) error {
		fmt.Printf("[%s] Transaction {%s} started!\n", startTimestamp.Format("2006-01-02 15:04:05"), transactionId)

		for _, entry := range transactionLoggerData.TransactionLogs {
			fmt.Printf("--> [%s] [%s] %s", entry.Timestamp.Format("2006-01-02 15:04:05"), entry.LoggerLevel, entry.Message)
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

		fmt.Printf("[%s] Transaction {%s} ended!\n", endTimeStamp.Format("2006-01-02 15:04:05"), transactionId)
		return nil
	}
}

func JSONTransactionLogOutputFileWrite() TransactionLogOutputWriter {
	return func(transactionId string, startTimeStamp time.Time, endTimeStamp time.Time, transactionLoggerData *TransactionLoggerData) error {
		// ! Not implemented
		return nil
	}
}

func TextTransactionLogOutputFileWrite() TransactionLogOutputWriter {
	return func(transactionId string, startTimeStamp time.Time, endTimeStamp time.Time, transactionLoggerData *TransactionLoggerData) error {
		// ! Not implemented
		return nil
	}
}
