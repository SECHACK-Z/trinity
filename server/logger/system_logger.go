package logger

import (
	"io"
	"log"
	"main/pubsub"
	"os"
)

type systemLogger struct {
	logFile *os.File
}

func (logger *systemLogger) fileLogger(event pubsub.SystemEvent) {
	log.Println(event.Type)
}

func StartSystemLogger() {
	systemLogFile, err := os.OpenFile("./systemLog", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		//log.Fatalf("cannot open systemlog")
	}
	logger := &systemLogger {
		logFile:systemLogFile,
	}
	log.SetOutput(io.MultiWriter(logger.logFile, os.Stdout))

	ps := pubsub.GetSystemEventPubSub()
	ps.Sub(logger.fileLogger)
}
