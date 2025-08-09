package provider

import (
	"fmt"
	"sms-dispatcher/logger"
	"sms-dispatcher/model"
)

func SendToMCI(msg model.SmsMessage) error {
	logger.Log.Infow("Sending to MCI", "phone", msg.Phone, "text", msg.Text)
	fmt.Println("[MCI] sent:", msg.Phone)
	return nil
}
