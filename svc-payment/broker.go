package main

import (
	"context"
	"log"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type broker struct {
	addr string
	conn *amqp091.Connection
	ch   *amqp091.Channel
}

func NewBroker(addr string) *broker {
	conn, err := amqp091.Dial(addr)
	if err != nil {
		log.Fatal(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}

	return &broker{
		addr: addr,
		conn: conn,
		ch:   ch,
	}
}

func (b *broker) Send(payload []byte) error {
	defer b.conn.Close()

	ctx, cancel := context.WithTimeout(context.TODO(), 3*time.Second)
	defer cancel()

	err := b.ch.PublishWithContext(ctx,
		"event",
		"event.*",
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        payload,
		},
	)

	if err != nil {
		log.Println(err)
		return nil
	}

	log.Println("sended successfully...")

	return nil
}
