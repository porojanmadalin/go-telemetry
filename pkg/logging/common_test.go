package logging

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"runtime"
)

func captureOutput(f func() error) (string, error) {
	orig := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	err := f()
	os.Stdout = orig
	w.Close()
	out, _ := io.ReadAll(r)
	return string(out), err
}

func getFunctionName(function interface{}) (string, error) {
	pc := reflect.ValueOf(function).Pointer()
	fn := runtime.FuncForPC(pc)
	if fn != nil {
		return fn.Name(), nil
	}
	return "", fmt.Errorf("error: could not get the function name")
}
