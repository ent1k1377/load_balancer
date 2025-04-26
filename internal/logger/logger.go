package logger

import (
	"log"
	"os"
)

var (
	debugLogger = log.New(os.Stdout, "[DEBUG] ", log.Ldate|log.Ltime|log.Lshortfile)
	infoLogger  = log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
	warnLogger  = log.New(os.Stderr, "[WARN] ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger = log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
	fatalLogger = log.New(os.Stderr, "[FATAL] ", log.Ldate|log.Ltime|log.Lshortfile)
)

func Debug(v ...any) {
	debugLogger.Println(v...)
}

func Debugf(format string, v ...any) {
	debugLogger.Printf(format, v...)
}

func Info(v ...any) {
	infoLogger.Println(v...)
}

func Infof(format string, v ...any) {
	infoLogger.Printf(format, v...)
}

func Warn(v ...any) {
	warnLogger.Println(v...)
}

func Warnf(format string, v ...any) {
	warnLogger.Printf(format, v...)
}

func Error(v ...any) {
	errorLogger.Println(v...)
}

func Errorf(format string, v ...any) {
	errorLogger.Printf(format, v...)
}

func Fatal(v ...any) {
	fatalLogger.Fatalln(v...)
}

func Fatalf(format string, v ...any) {
	fatalLogger.Fatalf(format, v...)
}
