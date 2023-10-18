package main

import (
	"fmt"
	"os"
	"time"
)

type Logger struct {
	file *os.File
}

func NewLogger(filename string) (*Logger, error) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if filename == "" { //if filename is empty, directly print to stdout
		file = os.Stdout
		err = nil
	}
	if err != nil {
		return nil, err
	}
	return &Logger{file}, nil
}

func (l *Logger) Info(message string, a ...any) {
	l.write("INFO", message, a...)
}
func (l *Logger) Error(message string, a ...any) {
	l.write("ERROR", message, a...)
}

func (l *Logger) write(level string, message string, a ...any) {
	timestamp := time.Now().Format("2006-01-31 15:04:05")
	message = fmt.Sprintf(message, a...)
	fmt.Fprintf(l.file, "[%s] %s: %s\n", timestamp, level, message)
}
