package api

import (
	"fmt"
	"net/http"
	"safelyyou_coding_challenge_go/database"
	"safelyyou_coding_challenge_go/models"
	"safelyyou_coding_challenge_go/services"

	"github.com/gin-gonic/gin"
)

// APIHandler holds the database connection and services, using dependency injection
type APIHandler struct {
	DB      *database.DatabaseOperations
	Metrics *services.MetricsService
}

// NewAPIHandler creates a new handler with the DB connection and services
func NewAPIHandler(db *database.DatabaseOperations, metrics *services.MetricsService) *APIHandler {
	return &APIHandler{
		DB:      db,
		Metrics: metrics,
	}
}

// Root handler for GET /
func (h *APIHandler) Root(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message":  "Fleet Monitoring API with SQLite",
		"version":  "1.0.0",
		"database": "SQLite",
	})
}

// ReceiveHeartbeat handler for POST /api/v1/devices/{device_id}/heartbeat
func (h *APIHandler) ReceiveHeartbeat(c *gin.Context) {
	deviceID := c.Param("device_id")

	if !h.DB.DeviceExists(deviceID) {
		c.JSON(http.StatusNotFound, gin.H{"detail": fmt.Sprintf("Device %s not found", deviceID)})
		return
	}

	var req models.HeartbeatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}

	if err := h.DB.AddHeartbeat(deviceID, req.SentAt); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to record heartbeat"})
		return
	}

	c.Status(http.StatusNoContent)
}

// ReceiveStats handler for POST /api/v1/devices/{device_id}/stats
func (h *APIHandler) ReceiveStats(c *gin.Context) {
	deviceID := c.Param("device_id")

	if !h.DB.DeviceExists(deviceID) {
		c.JSON(http.StatusNotFound, gin.H{"detail": fmt.Sprintf("Device %s not found", deviceID)})
		return
	}

	var req models.StatsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}

	if err := h.DB.AddStats(deviceID, req.SentAt, req.UploadTime); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to record stats"})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetDeviceStats handler for GET /api/v1/devices/{device_id}/stats
func (h *APIHandler) GetDeviceStats(c *gin.Context) {
	deviceID := c.Param("device_id")

	if !h.DB.DeviceExists(deviceID) {
		c.JSON(http.StatusNotFound, gin.H{"detail": fmt.Sprintf("Device %s not found", deviceID)})
		return
	}

	uptime := h.Metrics.CalculateUptime(deviceID)

	avgUploadTime := h.Metrics.CalculateAvgUploadTime(deviceID)

	c.JSON(http.StatusOK, models.StatsResponse{
		Uptime:        uptime,
		AvgUploadTime: avgUploadTime,
	})
}
