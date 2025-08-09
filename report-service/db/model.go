package db

type SmsStatus struct {
	ID        int64  `gorm:"primaryKey;autoIncrement"`
	MessageID string `json:"message_id" gorm:"index"`
	UserID    string `json:"user_id" gorm:"index"`
	Phone     string
	Status    string
	Timestamp int64
}
