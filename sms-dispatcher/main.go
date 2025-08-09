package main

import (
	"os"
	"os/signal"
	"sms-dispatcher/kafka"
	"sms-dispatcher/logger"
	"sms-dispatcher/provider"
	"syscall"
)

func main() {
	// 1. راه‌اندازی لاگر
	logger.SetupLogger()
	defer logger.Sync()

	// 2. بارگذاری فایل نگاشت اپراتورها
	err := provider.LoadOperatorConfig("config/operators.json")
	if err != nil {
		logger.Log.Fatalw("Failed to load operator config", "error", err)
	}

	// 3. Kafka broker ها
	brokers := []string{"localhost:9094"}

	// 4. راه‌اندازی کانسومرها
	go kafka.StartConsumer(brokers, "sms-normal")
	go kafka.StartConsumer(brokers, "sms-express")
	go kafka.StartDLQConsumer(brokers)

	logger.Log.Infow("sms-dispatcher service started")

	// 5. graceful shutdown
	waitForShutdown()
}

func waitForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	logger.Log.Infow("Waiting for shutdown signal...")
	<-quit
	logger.Log.Infow("sms-dispatcher shutting down gracefully...")
}
