package timer

import "time"

// GetNowTime GetNowTime
func GetNowTime() time.Time {
	location, _ := time.LoadLocation("Asia/Shanghai")
	return time.Now().In(location)
}

// GetCalculateTime GetCalculateTime
func GetCalculateTime(currentTime time.Time, d string) (time.Time, error) {
	duraction, err := time.ParseDuration(d)
	if err != nil {
		return time.Time{}, err
	}
	return currentTime.Add(duraction), nil
}
