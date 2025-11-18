package entity

import "time"

type EventLogEntity struct {
	ID         string    `bson:"id"`
	PlatformID string    `bson:"platform_id"`
	ProjectID  string    `bson:"project_id"`
	Token      string    `bson:"token"`
	EventType  string    `bson:"event_type"`
	CreatedAt  time.Time `bson:"created_at"`
}
