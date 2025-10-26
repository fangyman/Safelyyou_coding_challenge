package services

import (
	"math"
	"time"

	"safelyyou_coding_challenge_go/database"
)

// MetricsService handles device metrics calculations
type MetricsService struct {
	db *database.DatabaseOperations
}

// NewMetricsService creates a new metrics service
func NewMetricsService(db *database.DatabaseOperations) *MetricsService {
	return &MetricsService{db: db}
}

// CalculateUptime calculates device uptime percentage based on heartbeats
func (s *MetricsService) CalculateUptime(deviceID string) float64 {
	timestamps := s.db.GetHeartbeatTimestamps(deviceID)

	if len(timestamps) < 2 {
		return 0.0
	}

	firstHeartbeat := timestamps[0]
	lastHeartbeat := timestamps[len(timestamps)-1]

	totalMinutes := lastHeartbeat.Sub(firstHeartbeat).Minutes()

	if totalMinutes == 0 {
		return 100.0 // At least 2 heartbeats in 0 minutes is 100%
	}

	uptime := (float64(len(timestamps)) / totalMinutes) * 100.0
	return math.Min(uptime, 100.0)
}

// CalculateAvgUploadTime calculates average upload time for a device
func (s *MetricsService) CalculateAvgUploadTime(deviceID string) string {
	avgUploadTime := s.db.GetAverageUploadTime(deviceID)

	if !avgUploadTime.Valid {
		return "0s"
	}

	return time.Duration(avgUploadTime.Float64).String()
}
