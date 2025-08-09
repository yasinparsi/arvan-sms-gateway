package provider

import (
	"fmt"
	"sms-dispatcher/logger"
	"sms-dispatcher/model"
)

func SendToMTN(msg model.SmsMessage) error {
	logger.Log.Infow("Sending to MTN", "phone", msg.Phone, "text", msg.Text)
	fmt.Println("[MTN] sent:", msg.Phone)
	return nil
}
