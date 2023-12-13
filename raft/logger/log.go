package logger

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

var Glogger, _ = NewLogger("")

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

func (l *Logger) callerName() (name string) {
	pc, _, _, ok := runtime.Caller(3)
	if !ok {
		return "unknownFunc"
	}
	return runtime.FuncForPC(pc).Name()
}

func (l *Logger) write(level string, message string, a ...any) {
	name := l.callerName()
	name = strings.TrimPrefix(name, "lxtend.com/m/")
	timestamp := time.Now().Format("15:04:05")
	message = fmt.Sprintf(message, a...)
	fmt.Fprintf(l.file, "[%s] %s|%s: %s\n", timestamp, level, name, message)
}
