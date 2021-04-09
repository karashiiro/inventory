package main

import (
	log "github.com/sirupsen/logrus"
)

func main() {
	logFile := initLogging()
	defer logFile.Close()

	log.Info("Application started.")
}
