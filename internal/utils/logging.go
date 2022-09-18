package utils

import (
	"log"
	"os"
)

const (
	DebugLevel    = 10
	InfoLevel     = 20
	WarningLevel  = 30
	ErrorLevel    = 40
	CriticalLevel = 50
)

type Logger struct {
	level int

	debugLogger    *log.Logger
	infoLogger     *log.Logger
	warningLogger  *log.Logger
	errorLogger    *log.Logger
	criticalLogger *log.Logger
}

func (logger *Logger) SetLevel(level int) {
	logger.level = level
}

func (logger *Logger) Debug(message ...any) {
	if logger.level <= DebugLevel {
		logger.debugLogger.Println(message...)
	}
}

func (logger *Logger) Info(message ...any) {
	if logger.level <= InfoLevel {
		logger.infoLogger.Println(message...)
	}
}

func (logger *Logger) Warning(message ...any) {
	if logger.level <= WarningLevel {
		logger.warningLogger.Println(message...)
	}
}

func (logger *Logger) Error(message ...any) {
	if logger.level <= ErrorLevel {
		logger.errorLogger.Println(message...)
	}
}

func (logger *Logger) Critical(message ...any) {
	if logger.level <= CriticalLevel {
		logger.criticalLogger.Println(message...)
	}
}

func (logger *Logger) Initialize() {
	logger.debugLogger = log.New(os.Stdout, "[DEBUG] ", log.Ldate|log.Ltime)
	logger.infoLogger = log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime)
	logger.warningLogger = log.New(os.Stdout, "[WARNING] ", log.Ldate|log.Ltime)
	logger.errorLogger = log.New(os.Stdout, "[ERROR] ", log.Ldate|log.Ltime)
	logger.criticalLogger = log.New(
		os.Stdout,
		"[CRITICAL] ",
		log.Ldate|log.Ltime,
	)
}

func init() {
	Log.Initialize()
}

var Log Logger = Logger{}
