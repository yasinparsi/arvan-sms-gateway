package main

import (
	"log"
	"os"
	"strings"
	"sms-gateway/handler"
	"sms-gateway/kafka"
	"sms-gateway/sms-gateway/proto"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func main() {
	kafkaBrokersEnv := os.Getenv("KAFKA_BROKERS")
	var brokers []string
	if kafkaBrokersEnv == "" {
		brokers = []string{"Kafka00Service:9092"}
	} else {
		brokers = strings.Split(kafkaBrokersEnv, ",")
	}
	kafka.InitKafka(brokers)

	// connect to billing service using env-configured address
	billingAddr := os.Getenv("BILLING_GRPC_ADDR")
	if billingAddr == "" {
		billingAddr = "billing-service:4040"
	}
	conn, err := grpc.Dial(billingAddr, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(3*time.Second))
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
