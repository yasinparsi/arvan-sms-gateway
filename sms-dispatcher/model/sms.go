package model

type SmsMessage struct {
	UserID    string `json:"user_id"`
	Phone     string `json:"phone"`
	Text      string `json:"text"`
	Type      string `json:"type"` // normal or express
	MessageID string `json:"message_id"`
}
