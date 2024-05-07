package shared

import (
	"encoding/json"
	"fmt"
)

func JsonStringToObject(str string) map[string]interface{} {
	var result map[string]interface{}
	err := json.Unmarshal([]byte(str), &result)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return result
}

func ObjectToJsonString(obj map[string]interface{}) string {
	jsonData, err := json.Marshal(obj)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(jsonData)
}
