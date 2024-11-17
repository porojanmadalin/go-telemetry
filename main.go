package main

import (
	"fmt"
	"go-telemetry/pkg/logging"
)

func main() {

	// log := logging.NewLog()

	// log.Debug("Hello", map[string]any{
	// 	"VarInt":  1,
	// 	"VarStr":  "Test",
	// 	"VarTest": 3.14,
	// })
	// log.Warning("Hello", map[string]any{
	// 	"VarInt":  1,
	// 	"VarStr":  "Test",
	// 	"VarTest": 3.14,
	// })
	// log.Info("Hello", map[string]any{
	// 	"VarInt":  1,
	// 	"VarStr":  "Test",
	// 	"VarTest": 3.14,
	// })
	// log.Error("Hello", map[string]any{
	// 	"VarInt":  1,
	// 	"VarStr":  "Test",
	// 	"VarTest": 3.14,
	// })

	transactionLog := logging.NewTransactionLog("mainTest")
	err := transactionLog.StartTransactionLogging()
	if err != nil {
		fmt.Println(err)
	}

	transactionLog.Debug("Hello", map[string]any{
		"VarInt":  1,
		"VarStr":  "Test",
		"VarTest": 3.14,
	})
	transactionLog.Warning("Hello", map[string]any{
		"VarInt":  1,
		"VarStr":  "Test",
		"VarTest": 3.14,
	})
	transactionLog.Info("Hello", map[string]any{
		"VarInt":  1,
		"VarStr":  "Test",
		"VarTest": 3.14,
	})
	transactionLog.Error("Hello", map[string]any{
		"VarInt":  1,
		"VarStr":  "Test",
		"VarTest": 3.14,
	})
	err = transactionLog.StopTransactionLogging()
	if err != nil {
		fmt.Println(err)
	}
}
