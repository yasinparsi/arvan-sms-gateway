package main

import (
	"report-service/api"
	"report-service/db"
	"report-service/kafka"
	"report-service/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	logger.SetupLogger()
	defer logger.Sync()

	dsn := "postgres://report_user:report_pass@localhost:5433/sms_report?sslmode=disable"
	db.InitDB(dsn)

	go kafka.StartStatusConsumer([]string{"localhost:9094"})

	r := gin.Default()
	r.GET("/report/:user_id", api.GetSmsByUser)
	r.GET("/report/:user_id/export", api.ExportSmsCSV)

	r.Run(":6070")
}
