package main

import (
	"go-telemetry/pkg/logging"
)

func main() {
	//TODO: add a transaction started/finished fn

	log := logging.New()

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
}
