package viewmodel

import "time"

type ResourceChangeData struct {
	ID              string    `json:"id"`
	ResourceVersion int       `json:"version"`
	EventType       string    `json:"event_type"`
	CreatedAt       time.Time `json:"created_at"`
	Name            string    `json:"name"`
	Type            string    `json:"type"`
	Data            string    `json:"data"`
	ImageData       string    `json:"imageData"`
	Tags            []string  `json:"tags"`
}

type ResourceView struct {
	ID        string    `json:"id" redis:"id"`
	Name      string    `json:"name" redis:"name"`
	Type      string    `json:"type" redis:"type"`
	Data      string    `json:"data" redis:"data"`
	ImageData string    `json:"imageData" redis:"imageData"`
	Version   int       `json:"version" redis:"version"`
	IsDelete  bool      `json:"is_deleted" redis:"is_deleted"`
	CreatedAt time.Time `json:"created_at" redis:"created_at"`
	UpdatedAt time.Time `json:"updated_at" redis:"updated_at"`
	Tags      []string  `json:"tags" redis:"-"`
	TagString string    `json:"-" redis:"tag_string"`
}
