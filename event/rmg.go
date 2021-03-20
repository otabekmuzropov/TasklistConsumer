package event

import (
	"github.com/streadway/amqp"
	"log"
)

type RabbitMQ struct {
	Channel *amqp.Channel
	Connection *amqp.Connection
}

func NewRabbitMQ() RabbitMQ {
	conn, err := amqp.Dial("amqp://admin:admin@localhost:5672/")

	if err != nil {
		log.Panic("Failed to connection RabbitMQ", err)
	}

	ch, err := conn.Channel()

	if err != nil {
		log.Panic("Failed to create new channel", err)
	}

	err = ch.ExchangeDeclare(
		"course",
		amqp.ExchangeTopic,
		true,
		false,
		false,
		false,
		nil,
		)

	if err != nil {
		log.Panic("Error to declare exchange", err)
	}

	queue1, err := ch.QueueDeclare(
		"course.update",
		false,
		false,
		true,
		false,
		nil,
		)

	if err != nil {
		log.Panic("Error to declare queue", err)
	}

	err = ch.QueueBind(
		queue1.Name,
		"course.#",
		"course",
		false,
		nil,
		)

	if err != nil {
		log.Println("Error while binding queue", err)
	}

	queue2, err := ch.QueueDeclare(
		"course.update",
		false,
		false,
		true,
		false,
		nil,
		)

	if err != nil {
		log.Panic("Failed to declare to queue", err)
	}

	err = ch.QueueBind(
		queue2.Name,
		"course.update",
		"course",
		false,
		nil,
		)

	if err != nil {
		log.Panic(err)
	}
	return RabbitMQ{
		Channel:    ch,
		Connection: conn,
	}

}

func (r *RabbitMQ) Consume(q1, q2 string) {


	msgs1, err := r.Channel.Consume(
		q1,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Panic("Failed consuming")
	}

	msgs2, _ := r.Channel.Consume(
		q2,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Panic("Failed consuming")
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs1 {
			log.Printf(" [x] %s", d.Body)
		}
	}()

	go func() {
		for d := range msgs2 {
			log.Printf(" [x] %s", d.Body)
		}
	}()

	<-forever
}
