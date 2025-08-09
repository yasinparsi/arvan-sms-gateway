package main

import (
	"log"
	"sms-gateway/handler"
	"sms-gateway/kafka"
	"sms-gateway/sms-gateway/proto"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func main() {
	brokers := []string{"localhost:9094"}
	kafka.InitKafka(brokers)

	// اتصال به billing-service
	conn, err := grpc.Dial("localhost:4040", grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(3*time.Second))
	if err != nil {
		log.Fatalf("failed to connect to billing service: %v", err)
	}
	defer conn.Close()

	billingClient := proto.NewBillingServiceClient(conn)

	r := gin.Default()
	handler.SetBillingClient(billingClient)
	r.POST("/send", handler.SendSMS)

	r.Run(":7070")
}
