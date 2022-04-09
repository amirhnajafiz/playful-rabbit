package client

import (
	"fmt"
	"strings"

	"github.com/amirhnajafiz/playful-rabbit/internal/test"
	"github.com/streadway/amqp"
)

type Client struct {
	Cfg        Config
	Queue      string
	Connection *amqp.Connection
}

func (c *Client) Push(s string) error {
	ch, err := c.Connection.Channel()
	if err != nil {
		return err
	}

	defer func() {
		err := ch.Close()
		if err != nil {
			panic(err)
		}
	}()

	_, err = ch.QueueDeclare(
		c.Queue,
		c.Cfg.Durable,
		c.Cfg.AutoDelete,
		c.Cfg.Exclusive,
		!c.Cfg.Wait,
		nil,
	)
	if err != nil {
		return err
	}

	err = ch.Publish(
		"",
		c.Queue,
		c.Cfg.Mandatory,
		c.Cfg.Immediate,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(s),
		},
	)

	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Listen() error {
	ch, err := c.Connection.Channel()
	if err != nil {
		return err
	}

	defer func() {
		err := ch.Close()
		if err != nil {
			panic(err)
		}
	}()

	messages, err := ch.Consume(
		c.Queue,
		"",
		c.Cfg.AutoAck,
		c.Cfg.Exclusive,
		!c.Cfg.Local,
		!c.Cfg.Wait,
		nil,
	)

	if err != nil {
		return err
	}

	forever := make(chan bool)
	go func() {
		for d := range messages {
			parts := strings.Split(string(d.Body), " Brear ")

			test.Done(parts[0])

			fmt.Printf("[%s]: %s \n", parts[0], parts[1])
		}
	}()

	<-forever

	return nil
}
