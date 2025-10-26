package models

import "time"

type DeviceModel struct {
	DeviceID string `gorm:"primaryKey;index"`
}

type HeartbeatModel struct {
	DeviceID string `gorm:"index;"`
	ID       uint   `gorm:"primaryKey;autoincrement;"`
	SentAt   time.Time
}

type StatsModel struct {
	DeviceID   string    `gorm:"index;"`
	ID         uint      `gorm:"primaryKey;autoincrement;"`
	SentAt     time.Time `gorm:"column:sent_at"`
	UploadTime int64     `gorm:"column:upload_time"` // nanoseconds
}

func (DeviceModel) TableName() string    { return "devices" }
func (HeartbeatModel) TableName() string { return "heartbeats" }
func (StatsModel) TableName() string     { return "stats" }

// AllModels returns all database models for migration
func AllModels() []interface{} {
	return []interface{}{
		&DeviceModel{},
		&HeartbeatModel{},
		&StatsModel{},
	}
}

type HeartbeatRequest struct {
	SentAt time.Time `json:"sent_at"`
}

type StatsRequest struct {
	SentAt     time.Time `json:"sent_at"`
	UploadTime int64     `json:"upload_time" binding:"required"`
}

type StatsResponse struct {
	Uptime        float64 `json:"uptime"`
	AvgUploadTime string  `json:"avg_upload_time"`
}
