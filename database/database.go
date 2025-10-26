package database

import (
	"database/sql"
	"time"

	"safelyyou_coding_challenge_go/models"

	"gorm.io/gorm"
)

// DatabaseOperations holds the GORM database connection
type DatabaseOperations struct {
	DB *gorm.DB
}

// Device operations

// DeviceExists checks if a device is in the database
func (ops *DatabaseOperations) DeviceExists(deviceID string) bool {
	var device models.DeviceModel
	ops.DB.Model(&models.DeviceModel{}).Where("device_id = ?", deviceID).First(&device)
	return device.DeviceID == deviceID
}

// GetAllDeviceIDs returns all device IDs from the database
func (ops *DatabaseOperations) GetAllDeviceIDs() []string {
	var deviceIDs []string
	ops.DB.Model(&models.DeviceModel{}).Pluck("device_id", &deviceIDs)
	return deviceIDs
}

// CreateDevices bulk inserts devices into the database
func (ops *DatabaseOperations) CreateDevices(devices []models.DeviceModel) error {
	return ops.DB.CreateInBatches(&devices, 100).Error
}

// Heartbeat operations

// AddHeartbeat adds a new heartbeat record
func (ops *DatabaseOperations) AddHeartbeat(deviceID string, sentAt time.Time) error {
	heartbeat := models.HeartbeatModel{
		DeviceID: deviceID,
		SentAt:   sentAt,
	}
	return ops.DB.Create(&heartbeat).Error
}

// GetHeartbeatTimestamps returns all heartbeat timestamps for a device
func (ops *DatabaseOperations) GetHeartbeatTimestamps(deviceID string) []time.Time {
	var timestamps []time.Time
	ops.DB.Model(&models.HeartbeatModel{}).
		Where("device_id = ?", deviceID).
		Order("sent_at asc").
		Distinct().
		Pluck("sent_at", &timestamps)
	return timestamps
}

// Stats operations

// AddStats adds a new stats record
func (ops *DatabaseOperations) AddStats(deviceID string, sentAt time.Time, uploadTime int64) error {
	stats := models.StatsModel{
		DeviceID:   deviceID,
		SentAt:     sentAt,
		UploadTime: uploadTime,
	}
	return ops.DB.Create(&stats).Error
}

// GetAverageUploadTime returns the average upload time for a device
func (ops *DatabaseOperations) GetAverageUploadTime(deviceID string) sql.NullFloat64 {
	var avgUploadTime sql.NullFloat64
	ops.DB.Model(&models.StatsModel{}).
		Where("device_id = ?", deviceID).
		Select("AVG(upload_time) AS avg_upload_time").
		Scan(&avgUploadTime)
	return avgUploadTime
}
