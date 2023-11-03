package json

import (
	"encoding/json"
)

func ToJSONString(i interface{}) string {
	bytes, err := json.MarshalIndent(i, "", "  ")
	if err != nil {
		return err.Error()
	}
	return string(bytes)
}
