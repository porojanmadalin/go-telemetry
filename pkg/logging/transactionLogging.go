package logging

import (
	"fmt"
	"sync"
	"time"
)

type TransactionLoggerData struct {
	LoggerLevel     loggerLevel   `json:"loggerLevel"`
	TransactionLogs []*LoggerData `json:"transactionLogs"`
}

type transactionLogging struct {
	transactionId  string
	loggerLevel    loggerLevel
	startTimestamp time.Time
	outputWrite    TransactionLogOutputWriter
}

type transactionMap struct {
	sync.Map
}

var transactionLoggerOnce sync.Once
var transactionLoggerInstance *transactionLogging

var availableTransactions *transactionMap // using hash map for increased read/write performance

var addLogMutex sync.Mutex
var writeTransactionLogOutputMutex sync.Mutex

func NewTransactionLog(transactionId string, options ...func(*transactionLogging)) (*transactionLogging, error) {
	transactionLoggerOnce.Do(func() {
		initConfig()

		availableTransactions = &transactionMap{}

	})

	transactionLoggerInstance = &transactionLogging{}

	switch loggerConfig.Logger.Level {
	case string(LevelOff), string(LevelInfo), string(LevelWarning), string(LevelError), string(LevelDebug):
		transactionLoggerInstance.loggerLevel = loggerLevel(loggerConfig.Logger.Level)
	default:
		transactionLoggerInstance.loggerLevel = LevelInfo
	}

	switch loggerConfig.Logger.OutputWriter {
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

	if _, ok := availableTransactions.Load(transactionId); ok {
		// do not give out reference to the same transaction, on multiple calls
		return nil, fmt.Errorf("error: the provided transaction was already started %s", transactionId)
	}

	transactionLoggerInstance.transactionId = transactionId

	return transactionLoggerInstance, nil
}

func WithTransactionLoggerLevel(loggerLevel loggerLevel) func(*transactionLogging) {
	return func(l *transactionLogging) {
		l.loggerLevel = loggerLevel
	}
}

func WithTransactionLogOutputWriter(outputWriter TransactionLogOutputWriter) func(*transactionLogging) {
	return func(l *transactionLogging) {
		l.outputWrite = outputWriter
	}
}

func (l *transactionLogging) Info(msg string, v MetaData) {
	l.processLoggerData(LevelInfo, msg, v)
}

func (l *transactionLogging) Warning(msg string, v MetaData) {
	l.processLoggerData(LevelWarning, msg, v)
}

func (l *transactionLogging) Error(msg string, v MetaData) {
	l.processLoggerData(LevelError, msg, v)
}

func (l *transactionLogging) Debug(msg string, v MetaData) {
	l.processLoggerData(LevelDebug, msg, v)
}

func (l *transactionLogging) processLoggerData(loggerLevel loggerLevel, msg string, metaData MetaData) {
	if convertLoggerLevelToInt(loggerLevel) <= convertLoggerLevelToInt(l.loggerLevel) {
		err := l.addLogToTransaction(&LoggerData{LoggerLevel: loggerLevel, Timestamp: time.Now(), Message: msg, MetaData: metaData})
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (l *transactionLogging) addLogToTransaction(log *LoggerData) error {
	foundTransaction, ok := availableTransactions.Load(l.transactionId)
	if !ok {
		return fmt.Errorf("error: the provided transaction was not started or was recently ended %s", l.transactionId)
	}

	foundTransactionTyped, ok := foundTransaction.(*TransactionLoggerData)
	if !ok {
		return fmt.Errorf("error: the stored transaction is of unknown type %s", l.transactionId)
	}

	addLogMutex.Lock()
	newTransaction := &TransactionLoggerData{
		LoggerLevel:     foundTransactionTyped.LoggerLevel,
		TransactionLogs: append(foundTransactionTyped.TransactionLogs, log),
	}
	availableTransactions.CompareAndSwap(l.transactionId, foundTransactionTyped, newTransaction)
	addLogMutex.Unlock()
	return nil
}

func (l *transactionLogging) StartTransactionLogging() error {
	if l.loggerLevel == LevelOff {
		return nil
	}

	if _, loaded := availableTransactions.LoadOrStore(l.transactionId, &TransactionLoggerData{
		LoggerLevel:     l.loggerLevel,
		TransactionLogs: []*LoggerData{},
	}); loaded {
		return fmt.Errorf("error: the provided transaction was already started %s", l.transactionId)
	}

	l.startTimestamp = time.Now()

	return nil
}

func (l *transactionLogging) StopTransactionLogging() error {
	if l.loggerLevel == LevelOff {
		return nil
	}

	fmt.Printf("info: Will end logging transaction in 5 seconds %s\n", l.transactionId)
	time.Sleep(5 * time.Second)

	foundTransaction, loaded := availableTransactions.LoadAndDelete(l.transactionId)
	if !loaded {
		return fmt.Errorf("error: the provided transaction was not started or was recently ended %s", l.transactionId)
	}

	foundTransactionTyped, ok := foundTransaction.(*TransactionLoggerData)
	if !ok {
		return fmt.Errorf("error: the stored transaction is of unknown type %s", l.transactionId)
	}

	// found transaction should not contain logs that do not sattisfy the log level set prior

	writeTransactionLogOutputMutex.Lock()
	defer writeTransactionLogOutputMutex.Unlock()

	err := l.outputWrite(l.transactionId, l.startTimestamp, time.Now(), foundTransactionTyped)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
