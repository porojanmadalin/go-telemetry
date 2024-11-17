package main

import (
	"go-telemetry/pkg/logging"
)

func main() {

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
}
