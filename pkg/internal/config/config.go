package config

import (
	"log"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

const (
	configFilePathEnvKey = "GO_TELEMETRY_FILE_PATH"

	defaultConfigFileName = "telemetry-config.yml"
)

// A Logger is an environment values holder for logging
type Logger struct {
	Level        string `yaml:"level"`
	OutputWriter string `yaml:"outputWriter"`
	OutputDir    string `yaml:"outputDir"`
}

// A Config is a generic environment values holder
type Config struct {
	Logger Logger `yaml:"logger"`
}

var configOnce sync.Once
var LoggerConfig *Config

// Init uses singleton pattern in order to load the configuration from a YAML file
func Init() *Config {
	configOnce.Do(func() {
		LoggerConfig = loadConfig()
		if LoggerConfig == nil {
			LoggerConfig = &Config{}
		}
	})
	return LoggerConfig
}

// loadConfig loads a YAML file into memory, that contains the library configuration set by the user
func loadConfig() *Config {
	configFileName := os.Getenv(configFilePathEnvKey)
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
