package date

import (
	"testing"
	"time"
)

func TestTimeDuration(t *testing.T) {
	a := time.Now().Local().Format(time.RFC3339)
	time.Sleep(1 * time.Minute)

	laterTime, err := ParseTimestamp(a)
	if err != nil {
		t.Error(err.Error())
		return
	}

	duration := MinutesDifferenceFronNow(laterTime)
	if duration != 1 {
		t.Errorf("There supposed to be a 1 second delay, not %d", duration)
	}

	t.Log("Success")
}
