package kafkaext

import (
	"fmt"
	"github.com/segmentio/kafka-go"
	"net"
)

func EnsureKafkaTopic(brokerAddr, topic string, partitions int) error {
	conn, err := kafka.Dial("tcp", brokerAddr)
	if err != nil {
		return fmt.Errorf("dial kafka: %w", err)
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return fmt.Errorf("get controller: %w", err)
	}

	controllerConn, err := kafka.Dial("tcp", net.JoinHostPort(controller.Host, fmt.Sprintf("%d", controller.Port)))
	if err != nil {
		return fmt.Errorf("dial controller: %w", err)
	}
	defer controllerConn.Close()
	topicConfigs := []kafka.TopicConfig{{
		Topic:             topic,
		NumPartitions:     partitions,
		ReplicationFactor: 1,
	}}
	err = controllerConn.CreateTopics(topicConfigs...)
	if err != nil {
		return fmt.Errorf("create topic: %w", err)
	}

	return nil
}
