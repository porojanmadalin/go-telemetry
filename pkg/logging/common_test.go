package logging

import (
	"go-telemetry/pkg/internal/config"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func setupTestEnvironment(t *testing.T) {
	config.LoggerConfig = &config.Config{}
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("fatal: unable to identify current directory")
	}

	outputPath := filepath.Join(filepath.Dir(file), "../../test/")

	err := os.MkdirAll(outputPath, os.ModePerm)
	if err != nil {
		t.Fatalf("fatal: could not create test directory %v", err)
	}

	config.LoggerConfig.Logger.OutputDir = outputPath
}

func cleanup(t *testing.T, logFiles []string) {
	for _, logFile := range logFiles {
		err := os.Remove(filepath.Join(config.LoggerConfig.Logger.OutputDir, logFile))
		if err != nil {
			t.Errorf("error: could not delete test artifact test/%s %v", logFile, err)
		}
	}
}
