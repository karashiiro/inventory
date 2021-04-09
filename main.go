package main

import (
	"io"
	"os"

	log "github.com/sirupsen/logrus"
)

func initLogging() *os.File {
	log.SetReportCaller(true)
	log.SetFormatter(&log.JSONFormatter{})

	// Create log directory if it does not exist
	_, err := os.Stat("log/")
	if os.IsNotExist(err) {
		os.Mkdir("log/", 0644)
	}

	// Create and open log file
	var logFile *os.File
	logFile, err = os.OpenFile("log/inventory.log", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}

	log.SetOutput(io.MultiWriter(os.Stdout, logFile))

	return logFile
}

func main() {
	logFile := initLogging()
	defer logFile.Close()

	log.Info("Application started.")
}
