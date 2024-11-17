package logging

import "sync"

type TransactionLogOutputWriter func(*TransactionLoggerData) error

var writeToTransactionLogFileMutex sync.Mutex

func CLITransactionLogOutputWrite() TransactionLogOutputWriter {
	return func(transactionLoggerData *TransactionLoggerData) error {
		// ! Not implemented
		return nil
	}
}

func JSONTransactionLogOutputFileWrite() TransactionLogOutputWriter {
	return func(transactionLoggerData *TransactionLoggerData) error {
		writeToTransactionLogFileMutex.Lock()
		defer writeToTransactionLogFileMutex.Unlock()
		// ! Not implemented
		return nil
	}
}

func TextTransactionLogOutputFileWrite() TransactionLogOutputWriter {
	return func(transactionLoggerData *TransactionLoggerData) error {
		writeToTransactionLogFileMutex.Lock()
		defer writeToTransactionLogFileMutex.Unlock()
		// ! Not implemented
		return nil
	}
}
