package util

import (
	"encoding/json"
)

// ToJson 将对象序列化为 JSON 字符串
func ToJson[T any](v T) (string, error) {
	bytes, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// FromJson 将 JSON 字符串反序列化为对象
func FromJson[T any](jsonStr string) (T, error) {
	var result T
	err := json.Unmarshal([]byte(jsonStr), &result)
	return result, err
}
