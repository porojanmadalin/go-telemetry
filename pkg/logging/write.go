package logging

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	cli      OutputWriterType = "cli"
	jsonFile OutputWriterType = "jsonFile"
	textFile OutputWriterType = "textFile"
)

type OutputWriterType string

type OutputWriter func(*LoggerData) error

func CLIOutputWrite() OutputWriter {
	return func(loggerData *LoggerData) error {
		fmt.Printf("[%s] [%s] %s", loggerData.Timestamp.Format("2006-01-02 15:04:05.9999"), loggerData.LoggerLevel, loggerData.Message)
		for k, v := range loggerData.MetaData {
			typeName := reflect.TypeOf(v).Name()
			if strings.Contains(typeName, "int") {
				fmt.Printf(" [%s=%d]", k, v)
			} else if strings.Contains(typeName, "float") {
				fmt.Printf(" [%s=%f]", k, v)
			} else if typeName == "string" {
				fmt.Printf(" [%s=%s]", k, v)
			} else {
				fmt.Printf("\n")
				return fmt.Errorf("error: formatting for this variable type is not implement %s", typeName)
			}
		}
		fmt.Printf("\n")
		return nil
	}
}

func JSONOutputFileWrite() OutputWriter {
	return func(loggerData *LoggerData) error {
		return nil
	}
}

func TextOutputFileWrite() OutputWriter {
	return func(loggerData *LoggerData) error {
		return nil
	}
}
