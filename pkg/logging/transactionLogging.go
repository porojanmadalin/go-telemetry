package logging

type TransactionLoggerData struct {
	LoggerData
	TransactionID int `json:"transactionId"`
}
