package pipeline

import "time"

type Message struct {
	Id           string    `json:"id"`
	Content      string    `json:"content"`
	DequeueCount int       `json:"deliveries"`
	Created      time.Time `json:"created"`
	Schedule     time.Time `json:"schedule"`
	LockedUntil  time.Time `json:"lock"`
	LastError    string    `json:"error"`
}

func (m *Message) IsLocked() bool {
	return m.LockedUntil.After(time.Now())
}

func (m *Message) IsTimeToRun() bool {
	return m.Schedule.Before(time.Now())
}

func (m *Message) RunIn(delay time.Duration) {
	m.Schedule = time.Now().Add(delay)
}

func (m *Message) Lock(seconds int) {
	m.LockedUntil = time.Now().Add(time.Duration(seconds) * time.Second)
}

func (m *Message) Unlock() {
	m.LockedUntil = time.Time{}
}
