package logger

import (
	"fmt"
	"log"
	"os"
)

const defaultSettings = log.Ldate | log.Ltime | log.Lshortfile

var (
	debugLogger = log.New(os.Stdout, "[DEBUG] ", defaultSettings)
	infoLogger  = log.New(os.Stdout, "[INFO] ", defaultSettings)
	warnLogger  = log.New(os.Stderr, "[WARN] ", defaultSettings)
	errorLogger = log.New(os.Stderr, "[ERROR] ", defaultSettings)
	fatalLogger = log.New(os.Stderr, "[FATAL] ", defaultSettings)
)

func Debug(v ...any) {
	debugLogger.Output(2, fmt.Sprintln(v...))
}

func Debugf(format string, v ...any) {
	debugLogger.Output(2, fmt.Sprintf(format, v...)+"\n")
}

func Info(v ...any) {
	infoLogger.Output(2, fmt.Sprintln(v...))
}

func Infof(format string, v ...any) {
	infoLogger.Output(2, fmt.Sprintf(format, v...)+"\n")
}

func Warn(v ...any) {
	warnLogger.Output(2, fmt.Sprintln(v...))
}

func Warnf(format string, v ...any) {
	warnLogger.Output(2, fmt.Sprintf(format, v...)+"\n")
}

func Error(v ...any) {
	errorLogger.Output(2, fmt.Sprintln(v...))
}

func Errorf(format string, v ...any) {
	errorLogger.Output(2, fmt.Sprintf(format, v...)+"\n")
}

func Fatal(v ...any) {
	fatalLogger.Output(2, fmt.Sprintln(v...))
	os.Exit(1)
}

func Fatalf(format string, v ...any) {
	fatalLogger.Output(2, fmt.Sprintf(format, v...)+"\n")
	os.Exit(1)
}
