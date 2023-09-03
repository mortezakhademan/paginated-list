package dataType

import "time"

func ParseDateTime(datetime string) time.Time {
	t, err := time.Parse(time.RFC3339Nano, datetime)
	if err != nil {
		t, err = time.Parse("2006-1-2 15:4", datetime)
	}
	if err != nil {
		t, err = time.ParseInLocation("2006-1-2", datetime, time.Local)
	}
	return t
}
