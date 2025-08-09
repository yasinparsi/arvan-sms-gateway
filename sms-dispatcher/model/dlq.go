package model

type DLQMessage struct {
	SmsMessage
	RetryCount int    `json:"retry_count"`
	FailReason string `json:"fail_reason"`
}
