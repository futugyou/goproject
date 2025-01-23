package extensions

import (
	"errors"
	"fmt"
	"regexp"
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

func Int64ToTime(ts int64) time.Time {
	switch {
	case ts > 1e18: // 纳秒级
		return time.Unix(0, ts)
	case ts > 1e15: // 微秒级
		return time.Unix(0, ts*int64(time.Microsecond))
	case ts > 1e12: // 毫秒级
		return time.Unix(0, ts*int64(time.Millisecond))
	default: // 秒级
		return time.Unix(ts, 0)
	}
}

var timeFormats = []string{
	"2006-01-02T15:04:05Z07:00", // ISO 8601
	"2006-01-02 15:04:05",
	"2006-01-02 15:04",
	"2006-01-02 15",
	"2006-01-02",
	"2006/01/02 15:04:05",
	"2006/01/02 15:04",
	"2006/01/02 15",
	"2006/01/02",
	"2006.01.02 15:04:05",
	"2006.01.02 15:04",
	"2006.01.02 15",
	"2006.01.02",
}

var timeRegexps = []struct {
	regex   *regexp.Regexp
	formats []string
}{
	{
		regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),
		[]string{"2006-01-02T15:04:05Z07:00"},
	},
	{
		regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}$`),
		[]string{"2006-01-02 15:04:05"},
	},
	{
		regexp.MustCompile(`^\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2}$`),
		[]string{"2006/01/02 15:04:05"},
	},
	{
		regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`),
		[]string{"2006-01-02"},
	},
	{
		regexp.MustCompile(`^\d{4}/\d{2}/\d{2}$`),
		[]string{"2006/01/02"},
	},
	{
		regexp.MustCompile(`^\d{4}\.\d{2}\.\d{2}$`),
		[]string{"2006.01.02"},
	},
}

func getCurrentYearPrefix() int {
	currentYear := time.Now().Year()
	return currentYear / 100
}

func expandYear(str string) string {
	if len(str) == 2 {
		// Why is the number 20/19 fixed at the beginning? Because I won’t live to the year 2100.
		if str >= fmt.Sprintf("%02d", getCurrentYearPrefix()+50) {
			return "19" + str
		}
		return "20" + str
	}
	return str
}

func TroublesomeTimeParse(str string) (time.Time, error) {
	str = expandYear(str)

	for _, pattern := range timeRegexps {
		if pattern.regex.MatchString(str) {
			for _, layout := range pattern.formats {
				parsedTime, err := time.Parse(layout, str)
				if err == nil {
					return parsedTime, nil
				}
			}
		}
	}

	var lastErr error
	for _, layout := range timeFormats {
		parsedTime, err := time.Parse(layout, str)
		if err == nil {
			return parsedTime, nil
		}
		lastErr = err
	}

	return time.Time{}, errors.New("failed to parse time: " + str + " | last error: " + lastErr.Error())
}
