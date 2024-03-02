package date

import (
	"time"
)

func GenerateTimestamp() string {

	return time.Now().Format(time.RFC3339)
}

func ParseTimestamp(timestamp string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339, timestamp)

	if err != nil {
		return time.Time{}, err
	}

	return t.In(time.Local), nil
}

func SecondsDifferenceFromNow(a time.Time) int {
	return int(time.Since(a).Seconds())
}

func MinutesDifferenceFronNow(a time.Time) int {
	return int(time.Since(a).Minutes())
}

func HoursDifferenceFronNow(a time.Time) int {
	return int(time.Since(a).Hours())
}

func TimeDifference(a time.Time, b time.Time) time.Duration {
	diff := b.Sub(a)
	return diff
}
