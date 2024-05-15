package helpers

import (
	"encoding/json"
	"fmt"
)

func ToJson(data interface{}) []byte {
	v, err := json.Marshal(data)
	if err != nil {
		return nil
	}
	return v
}

func ToLogMessage(message string, data interface{}) string {
	return fmt.Sprintf("%s : %s", message, ToJsonString(data))
}

func ToJsonString(data interface{}) string {
	return string(ToJson(data))
}

func ToPrettyPrintString(data interface{}) string {
	v, err := json.MarshalIndent(data, "", "   ")
	if err != nil {
		return ""
	}
	return string(v)
}
