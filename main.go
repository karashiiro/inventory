package main

import (
	"encoding/json"
	"os"
	"strconv"
	"time"

	"github.com/karashiiro/inventory/message"
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

func logFailedAck(err error, corrID string) {
	log.Printf("message ack could not be delivered to channel\n\terror: %v\n\tcorrelation_id: ", err, corrID)
}

func parseRequestArgs(args []string) (uint32, uint32, error) {
	itemID, err := strconv.ParseUint(args[0], 10, 32)
	if err != nil {
		return 0, 0, err
	}

	quantity := uint64(0)
	if len(args) > 1 {
		quantity, err = strconv.ParseUint(args[1], 10, 32)
		if err != nil {
			return 0, 0, err
		}
	}

	return uint32(itemID), uint32(quantity), nil
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
			m := &message.Message{}
			err := json.Unmarshal(d.Body, m)
			if err != nil {
				log.Printf("failed to unmarshal message\n\terror: %v\n\tcorrelation_id: ", err, d.CorrelationId)
				err = d.Reject(false)
				if err != nil {
					logFailedAck(err, d.CorrelationId)
				}
				continue
			}

			switch m.Command {
			case "add":
				_, _, err = parseRequestArgs(m.Args)
				if err != nil {
					err = d.Reject(false)
					if err != nil {
						logFailedAck(err, d.CorrelationId)
					}
					continue
				}
			case "get":
				_, _, err = parseRequestArgs(m.Args)
				if err != nil {
					err = d.Reject(false)
					if err != nil {
						logFailedAck(err, d.CorrelationId)
					}
					continue
				}
			case "remove":
				_, _, err = parseRequestArgs(m.Args)
				if err != nil {
					err = d.Reject(false)
					if err != nil {
						logFailedAck(err, d.CorrelationId)
					}
					continue
				}
			default:
				log.Printf("failed to unmarshal message\n\tunk_msg: %v\n\tcorrelation_id: ", string(d.Body), d.CorrelationId)
				err = d.Reject(false)
				if err != nil {
					logFailedAck(err, d.CorrelationId)
				}
				continue
			}

			err = d.Ack(false)
			if err != nil {
				logFailedAck(err, d.CorrelationId)
			}
		}
	}()

	log.Info("Application started.")

	<-forever
}
