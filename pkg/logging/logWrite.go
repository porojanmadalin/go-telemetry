package logging

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"sync"
)

const (
	// JSON File Writer Constants
	startArray      = "[\n"
	endArray        = "\n]"
	indent          = "  "
	objectDelimiter = ",\n" + indent
)

type LogOutputWriter func(*LoggerData) error

var writeToLogFileMutex sync.Mutex

func CLILogOutputWrite() LogOutputWriter {
	return func(loggerData *LoggerData) error {
		fmt.Printf("[%s] [%s] %s", loggerData.Timestamp.Format("2006-01-02 15:04:05"), loggerData.LoggerLevel, loggerData.Message)
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

func JSONLogOutputFileWrite() LogOutputWriter {
	// TODO: take into account that multiple writes can occur on the same file
	return func(loggerData *LoggerData) error {
		writeToLogFileMutex.Lock()
		defer writeToLogFileMutex.Unlock()

		fileName := fmt.Sprintf("%s.json", loggerData.Timestamp.Format("2006-01-02"))
		f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			return fmt.Errorf("error: could not open json file %v", err)
		}
		defer f.Close()

		ret, err := f.Seek(-int64(len([]byte(endArray))), io.SeekEnd)
		if err != nil {
			_, err := f.Write([]byte(startArray))
			if err != nil {
				return fmt.Errorf("error: could not write json file %v", err)
			}
			ret = int64(len([]byte(startArray)))
		}

		loggerDataBytes, err := json.MarshalIndent(loggerData, indent, indent)
		if err != nil {
			return fmt.Errorf("error: could not marshal logger data %v", err)
		}

		writtenDataBytes := []byte(objectDelimiter)
		if ret == int64(len([]byte(startArray))) {
			writtenDataBytes = []byte(indent)
		}
		writtenDataBytes = append(writtenDataBytes, loggerDataBytes...)

		_, err = f.WriteAt(writtenDataBytes, ret)
		if err != nil {
			return fmt.Errorf("error: could not write json file %v", err)
		}

		_, err = f.WriteAt([]byte(endArray), ret+int64(len(writtenDataBytes)))
		if err != nil {
			return fmt.Errorf("error: could not write json file %v", err)
		}

		return nil
	}
}

func TextLogOutputFileWrite() LogOutputWriter {
	return func(loggerData *LoggerData) error {
		writeToLogFileMutex.Lock()
		defer writeToLogFileMutex.Unlock()

		fileName := fmt.Sprintf("%s.log", loggerData.Timestamp.Format("2006-01-02"))
		f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			return fmt.Errorf("error: could not open json file %v", err)
		}
		defer f.Close()

		f.WriteString(fmt.Sprintf("[%s] [%s] %s", loggerData.Timestamp.Format("2006-01-02 15:04:05"), loggerData.LoggerLevel, loggerData.Message))
		for k, v := range loggerData.MetaData {
			typeName := reflect.TypeOf(v).Name()
			if strings.Contains(typeName, "int") {
				f.WriteString(fmt.Sprintf(" [%s=%d]", k, v))
			} else if strings.Contains(typeName, "float") {
				f.WriteString(fmt.Sprintf(" [%s=%f]", k, v))
			} else if typeName == "string" {
				f.WriteString(fmt.Sprintf(" [%s=%s]", k, v))
			} else {
				f.WriteString("\n")
				return fmt.Errorf("error: formatting for this variable type is not implement %s", typeName)
			}
		}
		f.WriteString("\n")
		return nil
	}
}
