package config

import (
	"go-telemetry/pkg/logging"
	"log"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

const (
	ConfigFileNameEnvKey  = "GO_TELEMETRY_FILE_NAME"
	DefaultConfigFileName = "telemetry-config.yml"
)

type Config struct {
	Logging struct {
		Level logging.LoggingLevel `yaml:"level"`
	} `yaml:"logging"`
}

var configOnce sync.Once
var LoggingConfig Config

func Init() Config {
	configOnce.Do(func() {
		LoggingConfig = loadConfig()
	})
	return LoggingConfig
}

func loadConfig() Config {
	configFileName := os.Getenv(ConfigFileNameEnvKey)
	if configFileName == "" {
		configFileName = DefaultConfigFileName
	}

	f, err := os.Open(configFileName)
	if err != nil {
		log.Fatal("error: the logging config file could not be opened. Check if the file exists or if it is corrupt", err)
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		log.Fatal("error: the logging config could not be decoded. Check if the file exists or if it is corrupt", err)
	}
	return cfg
}
