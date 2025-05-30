package internal

import "github.com/segmentio/kafka-go"

func NewKafkaProducer() *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP("kafka:9092"),
		Topic:    "events",
		Balancer: &kafka.LeastBytes{},
	}
}
