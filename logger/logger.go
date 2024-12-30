package logger

import (
	"log"
	"os"
)

type Logger struct {
	fd *os.File
}

var logger Logger

func InitializeLogger(path string) error {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	logger.fd = file
	log.SetOutput(logger.fd)
	return nil
}

func Log(message string) {
	log.Println(message)
}

func ShutdownLogger() {
	Log("shutting down server...\n\n")
	logger.fd.Close()
}
