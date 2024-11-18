package logging

import (
	"fmt"
	"go-telemetry/pkg/internal/config"
	"sync"
	"time"
)

const (
	waitDurationUntilTransactionStop = 2 * time.Second
)

// A TransactionLoggerData holds the transaction logs for a transaction
type TransactionLoggerData struct {
	LoggerLevel     loggerLevel   `json:"loggerLevel"`
	TransactionLogs []*LoggerData `json:"transactionLogs"`
}

// A transactionLogging holds the top-level configuration of the transaction logger.
// They can be configured via the YAML configuration file or by pre-defined "drivers" or self-created ones.
// The startTimestamp is set when the Transaction is started.
type transactionLogging struct {
	transactionId  string
	loggerLevel    loggerLevel
	startTimestamp time.Time
	outputWrite    TransactionLogOutputWriter
}

// A transactionMap is an active transaction holder for the transaction logger.
type transactionMap struct {
	sync.Map
}

var transactionLoggerOnce sync.Once
var transactionLoggerInstance *transactionLogging

var availableTransactions *transactionMap // using hash map for increased read/write performance

var addLogMutex sync.Mutex
var writeTransactionLogOutputMutex sync.Mutex

// NewLog creates a transaction logging instance, respective to the defined YAML configuration or by given "drivers" in form of options argument.
// The "drivers" have a higher priority than YAML configuration.
//
// The transaction log configuration is set every time a transaction logging instance is defined.
//
// If a transaction with the same name is already present, the returned instance is nil.
//
// Default:
// If no configuration nor drivers are specified, the log level is Info with output to CLI.
func NewTransactionLog(transactionId string, options ...func(*transactionLogging)) (*transactionLogging, error) {
	transactionLoggerOnce.Do(func() {
		config.Init()

		availableTransactions = &transactionMap{}
	})

	transactionLoggerInstance = &transactionLogging{}

	switch config.LoggerConfig.Logger.Level {
	case string(LevelOff), string(LevelInfo), string(LevelWarning), string(LevelError), string(LevelDebug):
		transactionLoggerInstance.loggerLevel = loggerLevel(config.LoggerConfig.Logger.Level)
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

	if _, ok := availableTransactions.Load(transactionId); ok {
		// do not give out reference to the same transaction, on multiple calls
		return nil, fmt.Errorf("error: the provided transaction was already started %s", transactionId)
	}

	transactionLoggerInstance.transactionId = transactionId

	return transactionLoggerInstance, nil
}

// WithTransactionLoggerLevel is a pre-defined "driver" that specifies the transaction log level used
func WithTransactionLoggerLevel(loggerLevel loggerLevel) func(*transactionLogging) {
	return func(l *transactionLogging) {
		l.loggerLevel = loggerLevel
	}
}

// WithTransactionLogOutputWriter is a pre-defined "driver" that specifies the transaction output writer used
func WithTransactionLogOutputWriter(outputWriter TransactionLogOutputWriter) func(*transactionLogging) {
	return func(l *transactionLogging) {
		l.outputWrite = outputWriter
	}
}

// Info registers a log of level Info, for a specific transaction, with a message and additional attributes.
//
// If no attributes, use nil as MetaData
func (l *transactionLogging) Info(msg string, v MetaData) {
	l.processLoggerData(LevelInfo, msg, v)
}

// Warning registers a log of level Warning, for a specific transaction, with a message and additional attributes.
//
// If no attributes, use nil as MetaData
func (l *transactionLogging) Warning(msg string, v MetaData) {
	l.processLoggerData(LevelWarning, msg, v)
}

// Error registers a log of level Error, for a specific transaction, with a message and additional attributes.
//
// If no attributes, use nil as MetaData
func (l *transactionLogging) Error(msg string, v MetaData) {
	l.processLoggerData(LevelError, msg, v)
}

// Debug registers a log of level Debug, for a specific transaction, with a message and additional attributes.
//
// If no attributes, use nil as MetaData
func (l *transactionLogging) Debug(msg string, v MetaData) {
	l.processLoggerData(LevelDebug, msg, v)
}

// processLoggerData processes a transaction log triage.
// Logs are added to transaction if the set transaction log level is higher or equal than the log method used (Info, Warning, Error, Debug).
//
// If LogLevel is Off, no logs are kept.
func (l *transactionLogging) processLoggerData(loggerLevel loggerLevel, msg string, metaData MetaData) {
	if convertLoggerLevelToInt(loggerLevel) <= convertLoggerLevelToInt(l.loggerLevel) {
		err := l.addLogToTransaction(&LoggerData{LoggerLevel: loggerLevel, Timestamp: time.Now(), Message: msg, MetaData: metaData})
		if err != nil {
			fmt.Println(err)
		}
	}
}

// addLogToTransaction adds log to a specific transaction
//
// addLogToTransaction is safe to call concurrently with other operations and will
// block until all other operations finish.
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

// StartTransactionLogging initiates the transaction logging process. It loads the specific transaction into a synchronised Hash Map.
//
// If transaction exists or started already, an error will be returned.
//
// StartTransactionLogging is safe to call concurrently with other operations and will
// block until all other operations finish.
//
// ! Call StopTransactionLogging when logging is finished.
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

// StopTransactionLogging stops the transaction logging process. It deletes the specific transaction from the synchronised Hash Map.
// After deletion, the transaction is written to the output using the specified Transaction OutputWriter.
// Will block until the writing is finished.
//
// StopTransactionLogging is safe to call concurrently with other operations and will
// block until all other operations finish.
//
// If transaction does not exists or started already, an error will be returned.
func (l *transactionLogging) StopTransactionLogging() error {
	if l.loggerLevel == LevelOff {
		return nil
	}
	endTimestamp := time.Now()

	fmt.Printf("info: Will end logging transaction in a couple of seconds %s\n", l.transactionId)
	time.Sleep(waitDurationUntilTransactionStop)

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

	err := l.outputWrite(l.transactionId, l.startTimestamp, endTimestamp, foundTransactionTyped)
	if err != nil {
		fmt.Println(err)
		writeTransactionLogOutputMutex.Unlock()
		return err
	}

	writeTransactionLogOutputMutex.Unlock()
	return nil
}
