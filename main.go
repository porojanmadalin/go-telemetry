package main

import (
	"fmt"
	"go-telemetry/config"
)

func main() {
	config.Init()

	fmt.Println(config.LoggingConfig.Logging.Level)
}
