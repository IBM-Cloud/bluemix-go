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

//MapToKeys ...
func MapToKeys(i interface{}) []string {
	if i != nil {
		m, ok := i.(map[string]interface{})
		if ok {
			keys := make([]string, 0, len(m))
			for k := range m {
				keys = append(keys, k)
			}
			return keys
		}
	}
	return nil

}
