package config

import (
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func setupConfigFile(t *testing.T, logLevel string, outputWriter string) error {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("fatal: unable to identify current directory")
	}

	err := os.MkdirAll(filepath.Join(filepath.Dir(file), "../../../test/"), os.ModePerm)
	if err != nil {
		t.Fatalf("fatal: could not create test directory %v", err)
	}

	os.Setenv(configFilePathEnvKey, filepath.Join(filepath.Dir(file), "../../../test/"+defaultConfigFileName))
	configFilePath := os.Getenv(configFilePathEnvKey)
	f, err := os.OpenFile(configFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		t.Fatalf("fatal: could not open yml file %v", err)
	}
	defer f.Close()

	enc := yaml.NewEncoder(f)

	err = enc.Encode(Config{
		Logger: Logger{
			Level:        logLevel,
			OutputWriter: outputWriter,
		},
	})
	if err != nil {
		t.Fatalf("fatal: encoding failed %v", err)
	}
	t.Cleanup(func() {
		cleanup(t)
	})

	return nil
}

func cleanup(t *testing.T) {
	err := os.Remove(os.Getenv(configFilePathEnvKey))
	if err != nil {
		t.Fatalf("fatal: could not delete the testing config file %v", err)
	}
}

func TestInit(t *testing.T) {
	configOnce = sync.Once{}
	Init()

	assert.Equal(t, "", LoggerConfig.Logger.Level)
	assert.Equal(t, "", LoggerConfig.Logger.OutputWriter)

	configOnce = sync.Once{}
	setupConfigFile(t, "info", "cli")
	Init()

	assert.Equal(t, "info", LoggerConfig.Logger.Level)
	assert.Equal(t, "cli", LoggerConfig.Logger.OutputWriter)
}
