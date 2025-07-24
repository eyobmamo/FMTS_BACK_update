package main

import (
	"FMTS/initiator"
	"FMTS/kafka"
)

func main() {
	kafka.Kafka_demo()
	initiator.Initiator()

}
