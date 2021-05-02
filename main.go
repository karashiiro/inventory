package main

import (
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

func openChannel(url string, channelName string) (<-chan amqp.Delivery, error) {
	var conn *amqp.Connection
	var err error
	for conn == nil {
		conn, err = amqp.Dial(url)
		if err != nil {
			log.Println("couldn't open RabbitMQ connection, retrying in 5 seconds")
			time.Sleep(5 * time.Second)
		}
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	defer ch.Close()

	mq, err := ch.QueueDeclare(channelName, false, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	err = ch.Qos(1, 0, false)
	if err != nil {
		return nil, err
	}

	msgs, err := ch.Consume(mq.Name, "", false, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	return msgs, nil
}

func main() {
	// Start logger
	logFile, err := initLogging()
	if err != nil {
		log.Fatalln(err)
	}
	defer logFile.Close()

	// Connect to database
	_, err = initDatabase()
	if err != nil {
		log.Fatalln(err)
	}

	// Connect to message broker
	msgs, err := openChannel(os.Getenv("INVENTORY_RMQ_CONNECTION_STRING"), os.Getenv("INVENTORY_RMQ_CHANNEL"))
	if err != nil {
		log.Fatalln(err)
	}

	// Start message loop
	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Println(d)
		}
	}()

	log.Info("Application started.")

	<-forever
}
