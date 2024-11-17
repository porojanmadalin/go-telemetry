package logging

import (
	"fmt"
	"go-telemetry/config"
	"sync"
	"time"
)

type TransactionLoggerData struct {
	LoggerLevel LoggerLevel `json:"loggerLevel"`
	Timestamps  []time.Time `json:"timestamps"`
	Messages    []string    `json:"messages"`
	MetaDatas   []MetaData  `json:"metaDatas"`
}

type TransactionLog struct {
	transactionId string
	loggerLevel   LoggerLevel
	outputWrite   TransactionLogOutputWriter
}

type TransactionMap struct {
	sync.Map
}

var transactionLoggerOnce sync.Once
var transactionLoggerInstance *TransactionLog

var availableTransactions *TransactionMap // using hash map for increased read/write performance

var addLogMutex sync.Mutex

func NewTransactionLog(transactionId string, options ...func(*TransactionLog)) *TransactionLog {
	transactionLoggerOnce.Do(func() {
		config.Init()

		transactionLoggerInstance = &TransactionLog{}

		availableTransactions = &TransactionMap{}

		switch config.LoggerConfig.Logger.Level {
		case string(LevelOff), string(LevelInfo), string(LevelWarning), string(LevelError), string(LevelDebug):
			transactionLoggerInstance.loggerLevel = LoggerLevel(config.LoggerConfig.Logger.Level)
		default:
			transactionLoggerInstance.loggerLevel = LevelInfo
		}

		switch config.LoggerConfig.Logger.OutputWriter {
		case string(cli):
			transactionLoggerInstance.outputWrite = CLITransactionLogOutputWrite()
		case string(jsonFile):
			transactionLoggerInstance.outputWrite = JSONTransactionLogOutputFileWrite()
		case string(textFile):
			transactionLoggerInstance.outputWrite = TextTransactionLogOutputFileWrite()
		default:
			transactionLoggerInstance.outputWrite = CLITransactionLogOutputWrite()
		}

		// Log options override the YAML file configuration
		for _, o := range options {
			o(transactionLoggerInstance)
		}
	})

	transactionLoggerInstance.transactionId = transactionId

	return transactionLoggerInstance
}

func WithTransactionLoggerLevel(loggerLevel LoggerLevel) func(*TransactionLog) {
	return func(l *TransactionLog) {
		l.loggerLevel = loggerLevel
	}
}

func WithTransactionLogOutputWriter(outputWriter TransactionLogOutputWriter) func(*TransactionLog) {
	return func(l *TransactionLog) {
		l.outputWrite = outputWriter
	}
}

func (l *TransactionLog) Info(msg string, v MetaData) {
	l.processLoggerData(LevelInfo, msg, v)
}

func (l *TransactionLog) Warning(msg string, v MetaData) {
	l.processLoggerData(LevelWarning, msg, v)
}

func (l *TransactionLog) Error(msg string, v MetaData) {
	l.processLoggerData(LevelError, msg, v)
}

func (l *TransactionLog) Debug(msg string, v MetaData) {
	l.processLoggerData(LevelDebug, msg, v)
}

func (l *TransactionLog) processLoggerData(loggerLevel LoggerLevel, msg string, metaData MetaData) {
	if convertLoggerLevelToInt(loggerLevel) <= convertLoggerLevelToInt(l.loggerLevel) {
		err := l.addLogToTransaction(&LoggerData{LoggerLevel: loggerLevel, Message: msg, MetaData: metaData})
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (l *TransactionLog) addLogToTransaction(log *LoggerData) error {
	foundTransaction, ok := availableTransactions.Load(l.transactionId)
	if !ok {
		return fmt.Errorf("error: the provided transaction was not started or was recently ended  %s", l.transactionId)
	}

	foundTransactionTyped, ok := foundTransaction.(TransactionLoggerData)
	if !ok {
		return fmt.Errorf("error: the stored transaction is of unknown type  %s", l.transactionId)
	}

	addLogMutex.Lock()
	newTransaction := &TransactionLoggerData{
		LoggerLevel: foundTransactionTyped.LoggerLevel,
		Timestamps:  append(foundTransactionTyped.Timestamps, log.Timestamp),
		Messages:    append(foundTransactionTyped.Messages, log.Message),
		MetaDatas:   append(foundTransactionTyped.MetaDatas, log.MetaData),
	}
	availableTransactions.CompareAndSwap(l.transactionId, foundTransactionTyped, newTransaction)
	addLogMutex.Unlock()
	return nil
}

func (l *TransactionLog) StartTransaction() error {
	if l.loggerLevel == LevelOff {
		return nil
	}

	if _, loaded := availableTransactions.LoadOrStore(l.transactionId, &TransactionLoggerData{
		Timestamps:  []time.Time{},
		LoggerLevel: l.loggerLevel,
		Messages:    []string{},
		MetaDatas:   []MetaData{},
	}); loaded {
		return fmt.Errorf("error: the provided transaction was already started %s", l.transactionId)
	}

	return nil
}

func (l *TransactionLog) StopTransaction() error {
	if l.loggerLevel == LevelOff {
		return nil
	}

	fmt.Printf("info: Will end logging transaction in 5 seconds %s\n", l.transactionId)
	time.Sleep(5 * time.Second)

	foundTransaction, loaded := availableTransactions.LoadAndDelete(l.transactionId)
	if !loaded {
		return fmt.Errorf("error: the provided transaction was not started or was recently ended  %s", l.transactionId)
	}

	foundTransactionTyped, ok := foundTransaction.(TransactionLoggerData)
	if !ok {
		return fmt.Errorf("error: the stored transaction is of unknown type  %s", l.transactionId)
	}

	// found transaction should not contain logs that do not sattisfy the log level set prior
	err := l.outputWrite(&foundTransactionTyped)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
