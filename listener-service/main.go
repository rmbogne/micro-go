package main

import (
	"fmt"
	"listener/event"
	"log"
	"math"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// TODO 1-try to connect to Rabbitmq
	fmt.Println("Connecting to RabbitMQ, please wait....")
	rabbitConn, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()

	// TODO 2-start listening for messages
	log.Println("Listening for and consuming RabbitMQ messages...") 

	// TODO 3-create consumer
	csumer, err := event.NewConsumer(rabbitConn)
	if err != nil {
		panic(err)
	}

	// TODO 4-watch the queue and consume the events
	err = csumer.Listen([]string{"log.INFO", "log.WARNING", "log.ERROR"})
	if err != nil {
		panic(err)
	}
}

func connect() (*amqp.Connection, error) {
	var counts int64
	var backoff = 1 * time.Second
	var connection *amqp.Connection

	// we keep trying until we get the connection established
	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			fmt.Println("RabbitMQ not yet ready, please wait ...")
			counts++
		} else {
			log.Println("Connected to RabbitMQ....")
			connection = c
			break
		}

		if counts > 5 {
			fmt.Println(err)
			return nil, err
		}

		backoff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("backing off ....")
		time.Sleep(backoff)
		continue
	}
	return connection, nil
}
