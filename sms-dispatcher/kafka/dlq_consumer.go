package kafka

import (
	"context"
	"encoding/json"
	"sms-dispatcher/logger"
	"sms-dispatcher/model"
	"sms-dispatcher/provider"

	"github.com/segmentio/kafka-go"
)

const maxRetry = 3

func StartDLQConsumer(brokers []string) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   "sms-dlq",
		GroupID: "dlq-group",
	})

	logger.Log.Infow("DLQ Consumer started")

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			logger.Log.Errorw("DLQ Kafka read error", "error", err)
			continue
		}

		var msg model.DLQMessage
		if err := json.Unmarshal(m.Value, &msg); err != nil {
			logger.Log.Errorw("DLQ unmarshal failed", "error", err)
			continue
		}

		if msg.RetryCount >= maxRetry {
			logger.Log.Warnw("Retry limit reached, publishing to status", "message_id", msg.MessageID)
			publishToStatus(msg.SmsMessage, "failed")
			continue
		}

		err = provider.RouteAndSend(msg.SmsMessage)
		if err != nil {
			msg.RetryCount++
			msg.FailReason = err.Error()
			logger.Log.Errorw("Retry failed, requeueing to DLQ", "retry", msg.RetryCount)
			requeueToDLQ(msg)
		} else {
			publishToStatus(msg.SmsMessage, "success")
		}
	}
}
