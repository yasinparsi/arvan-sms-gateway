package provider

import (
	"fmt"
	"sms-dispatcher/logger"
	"sms-dispatcher/model"
)

func SendToRightel(msg model.SmsMessage) error {
	logger.Log.Infow("Sending to Rightel", "phone", msg.Phone, "text", msg.Text)
	fmt.Println("[Rightel] sent:", msg.Phone)
	return nil
}
