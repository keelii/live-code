package main

import (
	"encoding/json"
	"log"
)

func ToJsonBytes[T any](value T) []byte {
	data, err := json.Marshal(value)
	if err != nil {
		log.Println("ToJsonString error:", err)
		return []byte("")
	}
	return data
}
func ToJsonString[T any](value T) string {
	data, err := json.Marshal(value)
	if err != nil {
		log.Println("ToJsonString error:", err)
		return ""
	}
	return string(data)
}

func FromJsonString[T any](data string) T {
	var t T
	err := json.Unmarshal([]byte(data), &t)
	if err != nil {
		log.Println("FromJsonString error:", err)
		return t
	}
	return t
}
func FromJsonBytes[T any](data []byte) T {
	var t T
	err := json.Unmarshal(data, &t)
	if err != nil {
		log.Println("FromJsonBytes error:", err)
		return t
	}
	return t
}
