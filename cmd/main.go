package main

import (
	"bitbucket.org/alien_soft/TaskListConsumerRMQ/event"
)

func main() {

	rmq := event.NewRabbitMQ()
	defer rmq.Connection.Close()
	defer rmq.Channel.Close()

	rmq.Consume("course.update", "course.create")

}
