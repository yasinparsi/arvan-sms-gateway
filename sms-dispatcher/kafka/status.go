package kafka

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"sms-dispatcher/logger"
	"sms-dispatcher/model"
	"time"

	"github.com/segmentio/kafka-go"
)

func publishToStatus(msg model.SmsMessage, status string) {
	statusEvent := model.SmsStatus{
		MessageID: msg.MessageID,
		UserID:    msg.UserID,
		Phone:     msg.Phone,
		Status:    status,
		Timestamp: time.Now().Unix(),
	}

	data, err := json.Marshal(statusEvent)
	if err != nil {
		logger.Log.Errorw("Failed to marshal sms status", "error", err)
		return
	}

	// Determine Kafka brokers from env
	brokersEnv := os.Getenv("KAFKA_BROKERS")
	var brokers []string
	if brokersEnv == "" {
		brokers = []string{"Kafka00Service:9092"}
	} else {
		brokers = strings.Split(brokersEnv, ",")
	}

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  brokers,
		Topic:    "sms-status",
		Balancer: &kafka.LeastBytes{},
	})

	defer writer.Close()

	err = writer.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(msg.MessageID),
		Value: data,
		Time:  time.Now(),
	})

	if err != nil {
		logger.Log.Errorw("Failed to publish status", "error", err)
	} else {
		logger.Log.Infow("Published SMS status", "message_id", msg.MessageID, "status", status)
	}
}
