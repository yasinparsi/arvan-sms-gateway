package model

type SmsRequest struct {
    Phone      string `json:"phone" binding:"required"`
    Text       string `json:"text" binding:"required"`
    Type       string `json:"type" binding:"required"` // "express" or "normal"
    UserID     string `json:"user_id" binding:"required"`
    MessageID  string `json:"message_id"`
}
