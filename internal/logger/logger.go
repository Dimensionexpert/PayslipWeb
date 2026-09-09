package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

type Logger struct {
	fileLogger     *log.Logger
	terminalLogger *log.Logger
}

func New(logDir string) (*Logger, func() error, error) {
	if err := os.MkdirAll(logDir, 0o755); err != nil {
		return nil, nil, err
	}

	filename := filepath.Join(
		logDir,
		"payslip_"+time.Now().Format("2006-01-02_150405")+".log",
	)

	file, err := os.Create(filename)
	if err != nil {
		return nil, nil, err
	}

	fileLogger := log.New(
		file,
		"",
		log.Ldate|log.Ltime,
	)

	terminalLogger := log.New(
		os.Stdout,
		"",
		0,
	)

	return &Logger{
		fileLogger:     fileLogger,
		terminalLogger: terminalLogger,
	}, file.Close, nil
}

func (l *Logger) Info(format string, args ...any) {
	message := formatMessage(format, args...)

	l.fileLogger.Printf("INFO  %s", message)
	l.terminalLogger.Println(message)
}

func (l *Logger) Error(format string, args ...any) {
	message := formatMessage(format, args...)

	l.fileLogger.Printf("ERROR %s", message)
	l.terminalLogger.Println(message)
}

func (l *Logger) FileOnlyInfo(format string, args ...any) {
	l.fileLogger.Printf("INFO  "+format, args...)
}

func (l *Logger) FileOnlyError(format string, args ...any) {
	l.fileLogger.Printf("ERROR "+format, args...)
}

func formatMessage(format string, args ...any) string {
	return fmt.Sprintf(format, args...)
}
