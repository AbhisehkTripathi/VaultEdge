package logger

import (
	"log"
	"os"
)

const maxLogSize = 10 * 1024 * 1024

func InitLogFile() *os.File {
	if err := os.MkdirAll("logs", 0755); err != nil {
		log.Printf("Error creating logs directory: %v", err)
		return nil
	}

	logPath := "logs/app.log"

	fileInfo, err := os.Stat(logPath)
	if err == nil && fileInfo.Size() > maxLogSize {
		if err := os.Truncate(logPath, 0); err != nil {
			log.Printf("Error truncating log file: %v", err)
			return nil
		}
	}

	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("Error opening log file: %v", err)
		return nil
	}

	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Println("Application started")

	return logFile
}
