package main

import (
	"os"
	"report-service/api"
	"report-service/db"
	"report-service/kafka"
	"report-service/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	logger.SetupLogger()
	defer logger.Sync()

	// Postgres connection via env to work inside docker network
	pgHost := os.Getenv("POSTGRES_HOST")
	if pgHost == "" {
		pgHost = "PostgresService"
	}
	pgPort := os.Getenv("POSTGRES_PORT")
	if pgPort == "" {
		pgPort = "5432"
	}
	pgUser := os.Getenv("POSTGRES_USER")
	if pgUser == "" {
		pgUser = "report_user"
	}
	pgPass := os.Getenv("POSTGRES_PASSWORD")
	if pgPass == "" {
		pgPass = "report_pass"
	}
	pgDB := os.Getenv("POSTGRES_DB")
	if pgDB == "" {
		pgDB = "sms_report"
	}

	dsn := "postgres://" + pgUser + ":" + pgPass + "@" + pgHost + ":" + pgPort + "/" + pgDB + "?sslmode=disable"
	db.InitDB(dsn)

	// Kafka consumer should use kafka service name
	kafkaBrokers := []string{"Kafka00Service:9092"}
	go kafka.StartStatusConsumer(kafkaBrokers)

	r := gin.Default()
	r.GET("/report/:user_id", api.GetSmsByUser)
	r.GET("/report/:user_id/export", api.ExportSmsCSV)

	r.Run(":6070")
}
