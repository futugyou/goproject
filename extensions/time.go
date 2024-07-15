package extensions

import (
	"fmt"
	"time"
)

var time_format = []string{"2006-01-02 15:04", "2006-01-02 15:04:05", "2006-01-02"}

func ParseTime(value string) (time.Time, error) {
	if len(value) == 0 {
		return time.Time{}, nil
	}

	for _, f := range time_format {
		d, err := time.Parse(f, value)
		if err == nil {
			return d, nil
		}
	}

	return time.Time{}, fmt.Errorf("time parse error, raw data is %s", value)
}

func ParseTimeWithFormat(value string, format string) (time.Time, error) {
	if len(value) == 0 {
		return time.Time{}, nil
	}

	return time.Parse(format, value)
}
