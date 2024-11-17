package logging

import (
	"log"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

const (
	configFileNameEnvKey  = "GO_TELEMETRY_FILE_NAME"
	defaultConfigFileName = "telemetry-config.yml"
)

type config struct {
	Logger struct {
		Level        string `yaml:"level"`
		OutputWriter string `yaml:"outputWriter"`
	} `yaml:"logger"`
}

var configOnce sync.Once
var loggerConfig config

func initConfig() config {
	configOnce.Do(func() {
		loggerConfig = loadConfig()
	})
	return loggerConfig
}

func loadConfig() config {
	configFileName := os.Getenv(configFileNameEnvKey)
	if configFileName == "" {
		configFileName = defaultConfigFileName
	}

	f, err := os.Open(configFileName)
	if err != nil {
		log.Fatal("error: the logger config file could not be opened. Check if the file exists or if it is corrupt", err)
	}
	defer f.Close()

	var cfg config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		log.Fatal("error: the logger config could not be decoded. Check if the file exists or if it is corrupt", err)
	}
	return cfg
}
