package model

type SmsStatus struct {
	MessageID string `json:"message_id"`
	UserID    string `json:"user_id"`
	Phone     string `json:"phone"`
	Status    string `json:"status"` // success / failed
	Timestamp int64  `json:"timestamp"`
}
