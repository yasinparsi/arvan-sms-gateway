package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

var WriterNormal *kafka.Writer
var WriterExpress *kafka.Writer

func InitKafka(brokers []string) {
	WriterNormal = &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    "sms-normal",
		Balancer: &kafka.Hash{},
	}

	WriterExpress = &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    "sms-express",
		Balancer: &kafka.Hash{},
	}
}

func ProduceSMS(topic string, key string, msg interface{}) error {
	bytes, err := json.Marshal(msg)
	if err != nil {
		log.Printf("marshal error: %v", err)
		return err
	}

	writer := WriterNormal
	if topic == "sms-express" {
		writer = WriterExpress
	}

	err = writer.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(key),
		Value: bytes,
	})

	if err != nil {
		log.Printf("kafka write error: %v", err)
	}

	return err
}
