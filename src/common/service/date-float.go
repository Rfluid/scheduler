package common_service

import "time"

// Converts date to float using unix.
func DateToFloat(date time.Time) float64 {
	return float64(date.UnixNano()) / float64(time.Second)
}

// Converts float to date using unix.
func FloatToDate(date float64) time.Time {
	seconds := int64(date)
	nanos := int64((date - float64(seconds)) * float64(time.Second))

	return time.Unix(seconds, nanos)
}
