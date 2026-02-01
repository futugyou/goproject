package session

import (
	"time"

	"gorm.io/datatypes"
)

type SessionModel struct {
	AppName   string            `gorm:"primaryKey;index:idx_app_user"`
	UserID    string            `gorm:"primaryKey;index:idx_app_user"`
	SessionID string            `gorm:"primaryKey"`
	State     datatypes.JSONMap `gorm:"type:jsonb;default:'{}'"`
	UpdatedAt time.Time
	Events    []EventModel `gorm:"foreignKey:AppName,UserID,SessionID;references:AppName,UserID,SessionID;constraint:OnDelete:CASCADE"`
}

func (SessionModel) TableName() string { return "google_adk_sessions" }

type EventModel struct {
	ID        uint           `gorm:"primaryKey"`
	AppName   string         `gorm:"index"`
	UserID    string         `gorm:"index"`
	SessionID string         `gorm:"index"`
	EventData datatypes.JSON `gorm:"type:jsonb"`
	Timestamp time.Time      `gorm:"index"`
}

func (EventModel) TableName() string { return "google_adk_events" }
