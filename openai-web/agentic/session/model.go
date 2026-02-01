package session

import (
	"encoding/json"
	"fmt"
	"iter"
	"time"

	"google.golang.org/adk/session"
	"gorm.io/datatypes"
)

var _ session.Session = &SessionModel{}
var _ session.Events = &EventModelList{}
var _ session.State = &StateWrapper{}

type SessionModel struct {
	SessionID     string            `gorm:"primaryKey;column:id;type:varchar(64)"`
	SessionName   string            `gorm:"column:name;type:text;not null"`
	SessionUserID string            `gorm:"column:user_id;type:varchar(64);not null;index"`
	SessionState  datatypes.JSONMap `gorm:"column:state;type:jsonb;not null;default:'{}'"`
	UpdatedAt     time.Time         `gorm:"column:updated_at;type:timestamptz;not null"`
	EventsList    []EventModel      `gorm:"foreignKey:SessionID;references:SessionID;constraint:OnDelete:CASCADE"`
}

func (SessionModel) TableName() string { return "google_adk_sessions" }

type EventModel struct {
	ID        string         `gorm:"primaryKey;column:id;type:varchar(64)"`
	SessionID string         `gorm:"column:session_id;type:varchar(64);not null;index"`
	EventData datatypes.JSON `gorm:"column:data;type:jsonb;not null"`
	Timestamp time.Time      `gorm:"column:timestamp;type:timestamptz;not null;index"`
}

func (EventModel) TableName() string { return "google_adk_events" }

type StateWrapper struct {
	data datatypes.JSONMap
}

func (w *StateWrapper) Set(key string, value any) error {
	if w.data == nil {
		return fmt.Errorf("state data is nil")
	}
	w.data[key] = value
	return nil
}

func (w *StateWrapper) Get(key string) (any, error) {
	if w.data == nil {
		return nil, fmt.Errorf("state data is nil")
	}
	val, ok := w.data[key]
	if !ok {
		return nil, fmt.Errorf("key %q not found", key)
	}
	return val, nil
}

func (w *StateWrapper) All() iter.Seq2[string, any] {
	return func(yield func(string, any) bool) {
		for k, v := range w.data {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (s *SessionModel) AppName() string {
	return s.SessionName
}

func (s *SessionModel) ID() string {
	return s.SessionID
}

func (s *SessionModel) UserID() string {
	return s.SessionUserID
}

func (s *SessionModel) LastUpdateTime() time.Time {
	return s.UpdatedAt
}

func (s *SessionModel) State() session.State {
	if s.SessionState == nil {
		s.SessionState = make(datatypes.JSONMap)
	}
	return &StateWrapper{data: s.SessionState}
}

func (s *SessionModel) Events() session.Events {
	return &EventModelList{events: s.EventsList}
}

type EventModelList struct {
	events []EventModel
}

func (el *EventModelList) Len() int {
	return len(el.events)
}

func (el *EventModelList) At(i int) *session.Event {
	if i < 0 || i >= len(el.events) {
		return nil
	}
	row := el.events[i]

	var evt session.Event
	if len(row.EventData) > 0 {
		_ = json.Unmarshal(row.EventData, &evt)
	}

	evt.Timestamp = row.Timestamp
	return &evt
}

func (el *EventModelList) All() iter.Seq[*session.Event] {
	return func(yield func(*session.Event) bool) {
		for i := 0; i < len(el.events); i++ {
			if !yield(el.At(i)) {
				return
			}
		}
	}
}
