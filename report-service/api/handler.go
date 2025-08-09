package api

import (
	"encoding/csv"
	"net/http"
	"report-service/db"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func GetSmsByUser(c *gin.Context) {
	userID := c.Param("user_id")

	// فیلترهای زمانی
	startStr := c.Query("start")
	endStr := c.Query("end")

	// صفحه‌بندی
	pageStr := c.DefaultQuery("page", "1")
	sizeStr := c.DefaultQuery("size", "10")

	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}
	size, _ := strconv.Atoi(sizeStr)
	if size < 1 {
		size = 10
	}

	offset := (page - 1) * size

	query := db.DB.Where("user_id = ?", userID)

	if startStr != "" {
		start, err := strconv.ParseInt(startStr, 10, 64)
		if err == nil {
			query = query.Where("timestamp >= ?", start)
		}
	}
	if endStr != "" {
		end, err := strconv.ParseInt(endStr, 10, 64)
		if err == nil {
			query = query.Where("timestamp <= ?", end)
		}
	}

	var results []db.SmsStatus
	if err := query.
		Order("timestamp desc").
		Limit(size).
		Offset(offset).
		Find(&results).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"page":     page,
		"size":     size,
		"data":     results,
		"count":    len(results),
		"has_more": len(results) == size,
	})
}

func ExportSmsCSV(c *gin.Context) {
	userID := c.Param("user_id")
	startStr := c.Query("start")
	endStr := c.Query("end")

	query := db.DB.Where("user_id = ?", userID)

	if startStr != "" {
		start, err := strconv.ParseInt(startStr, 10, 64)
		if err == nil {
			query = query.Where("timestamp >= ?", start)
		}
	}
	if endStr != "" {
		end, err := strconv.ParseInt(endStr, 10, 64)
		if err == nil {
			query = query.Where("timestamp <= ?", end)
		}
	}

	var records []db.SmsStatus
	if err := query.Order("timestamp desc").Find(&records).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB query failed"})
		return
	}

	// تنظیم هدر برای دانلود CSV
	c.Header("Content-Disposition", "attachment; filename=report.csv")
	c.Header("Content-Type", "text/csv")

	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	// نوشتن عنوان ستون‌ها
	writer.Write([]string{"MessageID", "UserID", "Phone", "Status", "Timestamp"})

	// نوشتن رکوردها
	for _, r := range records {
		timestamp := time.Unix(r.Timestamp, 0).Format(time.RFC3339)
		writer.Write([]string{
			r.MessageID,
			r.UserID,
			r.Phone,
			r.Status,
			timestamp,
		})
	}
}
