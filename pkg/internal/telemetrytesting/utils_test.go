package telemetrytesting

import (
	"fmt"
	"os"

	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCaptureOutput(t *testing.T) {
	originalStdout := os.Stdout
	b, err := CaptureOutput(func() error {
		fmt.Println("test")
		return nil
	})
	if err != nil {
		t.Fatalf("fatal: could not capture the output %v", err)
	}

	assert.Equal(t, originalStdout, os.Stdout)
	assert.Contains(t, b, "test")
}

func TestGetFunctionName(t *testing.T) {
	fn1 := GetFunctionName
	name, err := GetFunctionName(fn1)
	if err != nil {
		t.Fatalf("fatal: could not get the function name %v", err)
	}
	assert.Contains(t, name, "GetFunctionName")
}

func TestGetFunctionNameWithUnnamedFunctions(t *testing.T) {
	// unnamed fn
	fn3 := func() {}
	name, err := GetFunctionName(fn3)
	if err != nil {
		t.Fatalf("fatal: could not get the function name %v", err)
	}
	assert.Contains(t, name, "func1")
}

func TestAlmostEqual(t *testing.T) {
	assert.True(t, AlmostEqual(3.14, 3.14))
}
