package core

import "time"

type TimeProvider interface {
	Now() time.Time
	Since(t time.Time) time.Duration
	Sleep(d time.Duration)
}

type RealTimeProvider struct{}

func (RealTimeProvider) Now() time.Time {
	return time.Now()
}

func (RealTimeProvider) Since(t time.Time) time.Duration {
	return time.Since(t)
}

func (RealTimeProvider) Sleep(d time.Duration) {
	time.Sleep(d)
}
