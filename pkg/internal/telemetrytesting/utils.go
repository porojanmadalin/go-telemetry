package abcd

import (
	"fmt"
	"io"
	"math"
	"os"
	"reflect"
	"runtime"
)

const float64EqualityThreshold = 1e-9

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

func GetFunctionName(function interface{}) (string, error) {
	pc := reflect.ValueOf(function).Pointer()
	fn := runtime.FuncForPC(pc)
	if fn != nil {
		return fn.Name(), nil
	}
	return "", fmt.Errorf("error: could not get the function name")
}

func AlmostEqual(a, b float64) bool {
	return math.Abs(a-b) <= float64EqualityThreshold
}
