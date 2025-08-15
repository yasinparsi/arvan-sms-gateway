package api

import (
	"net/http"
	"strconv"

	"billing-service/internal/storage"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	redisClient *storage.RedisClient
}

func NewHandler(redisClient *storage.RedisClient) *Handler {
	return &Handler{redisClient: redisClient}
}

// POST /charge/:userid?amount=10
func (h *Handler) ChargeUser(c *gin.Context) {
	userID := c.Param("userid")
	amountStr := c.Query("amount")
	if amountStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "amount query parameter is required"})
		return
	}

	amountInt, err := strconv.Atoi(amountStr) // convert to int
	if err != nil || amountInt <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "amount must be a positive integer"})
		return
	}

	err = h.redisClient.ChargeUser(c.Request.Context(), userID, int64(amountInt))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to charge user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "charged successfully", "user_id": userID, "amount": amountInt})
}
