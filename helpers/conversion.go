package helpers

import "time"

// Bool returns a pointer to the bool value
func Bool(v bool) *bool {
	return &v
}

// Int returns a pointer to the int value
func Int(v int) *int {
	return &v
}

// Duration returns a pointer to the time.Duration
func Duration(v time.Duration) *time.Duration {
	return &v
}
