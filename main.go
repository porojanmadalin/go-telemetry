package main

import (
	"fmt"
	"go-telemetry/pkg/logging"
	"sync"
)

func main() {

	// Examples

	//! #1 Logging example
	log := logging.NewLog()

	log.Debug("Hello", map[string]any{
		"VarInt":  1,
		"VarStr":  "Test",
		"VarTest": 3.14,
	})
	log.Warning("Hello", map[string]any{
		"VarInt":  1,
		"VarStr":  "Test",
		"VarTest": 3.14,
	})
	log.Info("Hello", map[string]any{
		"VarInt":  1,
		"VarStr":  "Test",
		"VarTest": 3.14,
	})
	log.Error("Hello", map[string]any{
		"VarInt":  1,
		"VarStr":  "Test",
		"VarTest": 3.14,
	})

	//! #2 Transaction Logging Example
	var waitGroup sync.WaitGroup
	waitGroup.Add(3)

	go func() {
		defer waitGroup.Done()

		transactionLog, err := logging.NewTransactionLog("mainTest1")
		if err != nil {
			fmt.Println(err)
			return
		}

		err = transactionLog.StartTransactionLogging()
		if err != nil {
			fmt.Println(err)
			return
		}

		transactionLog.Debug("Hello from 1", map[string]any{
			"VarInt":  1,
			"VarStr":  "Test",
			"VarTest": 3.14,
		})
		transactionLog.Warning("Hello from 1", map[string]any{
			"VarInt":  1,
			"VarStr":  "Test",
			"VarTest": 3.14,
		})
		transactionLog.Info("Hello from 1", map[string]any{
			"VarInt":  1,
			"VarStr":  "Test",
			"VarTest": 3.14,
		})
		transactionLog.Error("Hello from 1", map[string]any{
			"VarInt":  1,
			"VarStr":  "Test",
			"VarTest": 3.14,
		})
		err = transactionLog.StopTransactionLogging()
		if err != nil {
			fmt.Println(err)
			return
		}
	}()

	go func() {
		defer waitGroup.Done()

		transactionLog, err := logging.NewTransactionLog("mainTest2")
		if err != nil {
			fmt.Println(err)
			return
		}

		err = transactionLog.StartTransactionLogging()
		if err != nil {
			fmt.Println(err)
			return
		}

		transactionLog.Debug("Hello from 2", map[string]any{
			"VarInt":  1,
			"VarStr":  "Test",
			"VarTest": 3.14,
		})
		transactionLog.Warning("Hello from 2", map[string]any{
			"VarInt":  1,
			"VarStr":  "Test",
			"VarTest": 3.14,
		})
		transactionLog.Info("Hello from 2", map[string]any{
			"VarInt":  1,
			"VarStr":  "Test",
			"VarTest": 3.14,
		})
		transactionLog.Error("Hello from 2", map[string]any{
			"VarInt":  1,
			"VarStr":  "Test",
			"VarTest": 3.14,
		})
		err = transactionLog.StopTransactionLogging()
		if err != nil {
			fmt.Println(err)
			return
		}
	}()

	go func() {
		defer waitGroup.Done()

		transactionLog, err := logging.NewTransactionLog("mainTest3")
		if err != nil {
			fmt.Println(err)
			return
		}

		err = transactionLog.StartTransactionLogging()
		if err != nil {
			fmt.Println(err)
			return
		}

		transactionLog.Debug("Hello from 3", map[string]any{
			"VarInt": 1,
			"VarStr": "Test",
		})
		transactionLog.Warning("Hello from 3", map[string]any{
			"VarStr": "Test",
		})
		transactionLog.Info("Hello from 3", map[string]any{
			"VarTest": 3.14,
		})
		transactionLog.Error("Hello from 3", map[string]any{
			"VarInt": 1,
		})
		err = transactionLog.StopTransactionLogging()
		if err != nil {
			fmt.Println(err)
			return
		}
	}()
	waitGroup.Wait()
}
