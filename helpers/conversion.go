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

// String returns a pointer to the string value
func String(v string) *string {
	return &v
}

// Duration returns a pointer to the time.Duration
func Duration(v time.Duration) *time.Duration {
	return &v
}

// Map returns a pointer to the map
func Map(v map[string]interface{}) *map[string]interface{} {
	return &v
}

// List returns a pointer to the list
func List(v []string) *[]string {
	return &v
}
