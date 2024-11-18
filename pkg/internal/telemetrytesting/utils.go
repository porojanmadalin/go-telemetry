package telemetrytesting

import (
	"fmt"
	"io"
	"math"
	"os"
	"reflect"
	"runtime"
)

const float64EqualityThreshold = 1e-9

// CaptureOutput redirects the stdout from CLI to a byte array
func CaptureOutput(f func() error) (string, error) {
	orig := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	err := f()
	os.Stdout = orig
	w.Close()
	out, _ := io.ReadAll(r)
	return string(out), err
}

// GetFunctionName returns the full name (including corresponding packages) of the given argument function
func GetFunctionName(function interface{}) (string, error) {
	pc := reflect.ValueOf(function).Pointer()
	fn := runtime.FuncForPC(pc)
	if fn != nil {
		return fn.Name(), nil
	}
	return "", fmt.Errorf("error: could not get the function name")
}

// AlmostEqual compares two floats, with an approximation error of 1e-9
func AlmostEqual(a, b float64) bool {
	return math.Abs(a-b) <= float64EqualityThreshold
}
