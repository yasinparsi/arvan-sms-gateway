package kafka

import (
	"context"
	"encoding/json"
	"sms-dispatcher/logger"
	"sms-dispatcher/model"
	"time"

	"github.com/segmentio/kafka-go"
)

func publishToDLQ(msg model.SmsMessage, reason string) {
	dlqMsg := model.DLQMessage{
		SmsMessage: msg,
		RetryCount: 1,
		FailReason: reason,
	}

	data, err := json.Marshal(dlqMsg)
	if err != nil {
		logger.Log.Errorw("Failed to marshal DLQ message", "error", err)
		return
	}

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{"localhost:9094"},
		Topic:    "sms-dlq",
		Balancer: &kafka.LeastBytes{},
	})

	defer writer.Close()

	err = writer.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(msg.MessageID),
		Value: data,
		Time:  time.Now(),
	})

	if err != nil {
		logger.Log.Errorw("Failed to write to DLQ", "error", err)
	} else {
		logger.Log.Infow("Published to DLQ", "message_id", msg.MessageID, "reason", reason)
	}
}

func requeueToDLQ(msg model.DLQMessage) {
	data, err := json.Marshal(msg)
	if err != nil {
		logger.Log.Errorw("Failed to marshal requeued DLQ message", "error", err)
		return
	}

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{"localhost:9094"},
		Topic:    "sms-dlq",
		Balancer: &kafka.LeastBytes{},
	})

	defer writer.Close()

	err = writer.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(msg.MessageID),
		Value: data,
		Time:  time.Now(),
	})

	if err != nil {
		logger.Log.Errorw("Failed to requeue DLQ message", "error", err)
	} else {
		logger.Log.Infow("Requeued to DLQ", "message_id", msg.MessageID, "retry", msg.RetryCount)
	}
}
