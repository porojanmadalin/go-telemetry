package config

import (
	"log"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

const (
	configFileNameEnvKey  = "GO_TELEMETRY_FILE_PATH"
	defaultConfigFileName = "telemetry-config.yml"
)

type Logger struct {
	Level        string `yaml:"level"`
	OutputWriter string `yaml:"outputWriter"`
}

type Config struct {
	Logger Logger `yaml:"logger"`
}

var configOnce sync.Once
var LoggerConfig *Config

func Init() *Config {
	configOnce.Do(func() {
		LoggerConfig = loadConfig()
		if LoggerConfig == nil {
			LoggerConfig = &Config{}
		}
	})
	return LoggerConfig
}

func loadConfig() *Config {
	configFileName := os.Getenv(configFileNameEnvKey)
	if configFileName == "" {
		configFileName = defaultConfigFileName
	}

	f, err := os.Open(configFileName)
	if err != nil {
		log.Println("warning: the logger config file could not be opened. Check if the file exists or if it is corrupt", err)
		return nil
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		log.Println("warning: the logger config could not be decoded. Check if the file exists or if it is corrupt", err)
		return nil
	}
	return &cfg
}
