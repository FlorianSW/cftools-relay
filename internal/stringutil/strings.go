package stringutil

import (
	"encoding/json"
	"fmt"
	"strconv"
)

func Itos(v interface{}) string {
	switch value := v.(type) {
	case float32:
	case float64:
		return fmt.Sprintf("%f", value)
	case int:
		return strconv.Itoa(value)
	case string:
		return value
	case json.Number:
		return value.String()
	}
	return ""
}

func Itof(v interface{}) float64 {
	switch value := v.(type) {
	case float64:
		return value
	case int:
		return float64(value)
	case float32:
		return float64(value)
	case string:
		r, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return -1
		}
		return r
	case json.Number:
		r, err := value.Float64()
		if err != nil {
			return -1
		}
		return r
	}
	return -1
}
