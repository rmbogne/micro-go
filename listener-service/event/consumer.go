package event

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn      *amqp.Connection
	queueName string
}

func NewConsumer(conn *amqp.Connection) (Consumer, error) {
	log.Println("Setting up a new consumer, please wait...")
	consumer := Consumer{
		conn: conn,
	}
	err := consumer.setup()
	if err != nil {
		log.Println("Issue with consumer setup")
		return Consumer{}, err
	}
	log.Println("Sucessfully setup a new consumer")
	return consumer, nil
}

func (Consumer *Consumer) setup() error {
	log.Println("Setting up a new channel, please wait...")
	channel, err := Consumer.conn.Channel()
	if err != nil {
		log.Println("Not able to connect to the channel")
		return err
	}
	log.Println("Sucessfully setup a new channel")
	return declareExchange(channel)
}

type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (consumer *Consumer) Listen(topics []string) error {
	log.Println("trying consumming from new channel")
	ch, err := consumer.conn.Channel()
	if err != nil {
		log.Println("error consumming from new channel")
		return err
	}
	defer ch.Close()

	log.Println("Declaring a queue, please wait...")
	q, err := declareRandomQueue(ch)
	if err != nil {
		log.Println("Error declaring queue.")
		return err
	}

	log.Println("Binding topics to queue , please wait...")
	for _, s := range topics {
		ch.QueueBind(
			q.Name,
			s,
			"logs_topic",
			false,
			nil,
		)
		/*if err != nil {
			log.Println("Error Binding topics to queue ")
			return err
		}*/
	}
	messages, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return err
	}

	forever := make(chan bool)
	go func() {
		log.Println("starting consuming messages ")
		for d := range messages {
			var payload Payload
			_ = json.Unmarshal(d.Body, &payload)
			log.Println(d.Body)

			go handlePayload(payload)
		}
	}()
	fmt.Printf("waiting for message[Exchange, Queue] [logs_topic, %s]\n", q.Name)
	<-forever

	return nil
}

func handlePayload(payload Payload) {
	switch payload.Name {
	case "log", "event":
		//log whatever we get
		log.Println("Entering handle payload ")
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}
	case "auth":
		// authenticate here

	case "mail":

	default:
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}

	}
}

func logEvent(entry Payload) error {
	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	logServiceURL := "http://logger-service/log"

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	response, err := client.Do(request)

	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		return err
	}

	return nil

}
