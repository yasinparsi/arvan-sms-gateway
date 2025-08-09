package kafka

import (
	"context"
	"encoding/json"
	"report-service/db"
	"report-service/logger"

	"github.com/segmentio/kafka-go"
)

func StartStatusConsumer(brokers []string) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   "sms-status",
		GroupID: "report-service-group",
	})

	logger.Log.Infow("Started Kafka consumer for sms-status")

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			logger.Log.Errorw("Kafka read error", "error", err)
			continue
		}

		var status db.SmsStatus
		if err := json.Unmarshal(m.Value, &status); err != nil {
			logger.Log.Errorw("Unmarshal error", "error", err)
			continue
		}

		if err := db.DB.Create(&status).Error; err != nil {
			logger.Log.Errorw("DB insert failed", "error", err)
		} else {
			logger.Log.Infow("Saved sms status", "message_id", status.MessageID)
		}
	}
}
