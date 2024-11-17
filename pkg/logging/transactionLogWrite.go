package logging

type TransactionLogOutputWriter func(*TransactionLoggerData) error

func CLITransactionLogOutputWrite() TransactionLogOutputWriter {
	return func(loggerData *TransactionLoggerData) error {
		// ! Not implemented
		return nil
	}
}

func JSONTransactionLogOutputFileWrite() TransactionLogOutputWriter {
	return func(loggerData *TransactionLoggerData) error {
		// ! Not implemented
		return nil
	}
}

func TextTransactionLogOutputFileWrite() TransactionLogOutputWriter {
	return func(loggerData *TransactionLoggerData) error {
		// ! Not implemented
		return nil
	}
}
