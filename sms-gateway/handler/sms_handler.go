package handler

import (
	"context"
	"net/http"
	"sms-gateway/kafka"
	"sms-gateway/model"

	"sms-gateway/sms-gateway/proto"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var billingClient proto.BillingServiceClient

// SetBillingClient sets the gRPC billing client
func SetBillingClient(client proto.BillingServiceClient) {
	billingClient = client
}

func SendSMS(c *gin.Context) {
	var userReq model.SmsRequest

	if err := c.ShouldBindJSON(&userReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cost := int64(1)
	res, err := billingClient.CheckBalance(context.Background(), &proto.BillingRequest{
		UserId: userReq.UserID,
		Cost:   cost,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Billing error"})
		return
	}
	if !res.Allowed {
		c.JSON(http.StatusPaymentRequired, gin.H{"error": "Insufficient balance"})
		return
	}
	if userReq.MessageID == "" {
		userReq.MessageID = uuid.New().String()
	}

	// انتخاب topic
	topic := "sms-normal"
	if userReq.Type == "express" {
		topic = "sms-express"
	}

	// ارسال به Kafka
	err = kafka.ProduceSMS(topic, userReq.UserID, userReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to enqueue SMS"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":     "queued",
		"message_id": userReq.MessageID,
	})
}
