package kafka

import (
	"context"
	"encoding/json"
	"sms-dispatcher/logger"
	"sms-dispatcher/model"
	"sms-dispatcher/provider"

	"github.com/segmentio/kafka-go"
)

func StartConsumer(brokers []string, topic string) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  topic + "-group",
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})

	logger.Log.Infow("Started Kafka consumer", "topic", topic)

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			logger.Log.Errorw("Kafka read error", "topic", topic, "error", err)
			continue
		}

		var msg model.SmsMessage
		if err := json.Unmarshal(m.Value, &msg); err != nil {
			logger.Log.Errorw("Failed to unmarshal message", "error", err)
			continue
		}

		err = provider.RouteAndSend(msg)
		if err != nil {
			logger.Log.Errorw("Provider send failed", "phone", msg.Phone, "error", err)

			publishToDLQ(msg, err.Error()) // بیاریمش توی dlq
		} else {
			publishToStatus(msg, "success")
		}
	}
}
