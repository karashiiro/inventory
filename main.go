package main

import (
	log "github.com/sirupsen/logrus"
)

func main() {
	logFile, err := initLogging()
	if err != nil {
		log.Fatalln(err)
	}
	defer logFile.Close()

	_, err = initDatabase()
	if err != nil {
		log.Fatalln(err)
	}

	log.Info("Application started.")
}
