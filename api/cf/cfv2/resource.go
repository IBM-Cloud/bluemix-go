package cfv2

import "encoding/json"

//Metadata ...
type Metadata struct {
	GUID string `json:"guid"`
	URL  string `json:"url"`
}

//Resource ...
type Resource struct {
	Metadata Metadata
}

//NumberToInt64 ...
func NumberToInt64(number json.Number, defaultValue int64) int64 {
	if number != "" {
		i, err := number.Int64()
		if err == nil {
			return i
		}
	}
	return defaultValue
}

//NumberToInt ...
func NumberToInt(number json.Number, defaultValue int) int {
	if number != "" {
		i, err := number.Int64()
		if err == nil {
			return int(i)
		}
	}
	return defaultValue
}
