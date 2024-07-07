package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

func main() {
	addr := "amqp://testing:secrettesting@localhost:5672/"

	conn, err := amqp091.Dial(addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	dlv, err := ch.Consume(
		"notify",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	for d := range dlv {
		var body map[string]interface{}

		err := json.Unmarshal(d.Body, &body)
		if err != nil {
			log.Println(err)
		}

		fmt.Println(body)
	}

}
