package rgitf

import (
	"errors"
	"strconv"
)

func InterfaceToInt(param interface{}) (int64, error) {
	var thisInt int64
	switch param.(type) {
	case string:
		if value, ok := param.(string); ok {
			strToInt, _ := strconv.Atoi(value)
			thisInt = int64(strToInt)
		}
	case float64:
		if value, ok := param.(float64); ok {
			thisInt = int64(value)
		}
	case int:
		if value, ok := param.(int); ok {
			thisInt = int64(value)
		}
	case int64:
		if value, ok := param.(int64); ok {
			thisInt = value
		}
	case int32:
		if value, ok := param.(int32); ok {
			thisInt = int64(value)
		}
	case float32:
		if value, ok := param.(float32); ok {
			thisInt = int64(value)
		}
	case bool:
		if value, ok := param.(bool); ok {
			if value == true {
				thisInt = 1
			} else {
				thisInt = 0
			}
		}
	default:
		return 0, errors.New("unknow type")
	}
	return thisInt, nil
}

func InterfaceToString(param interface{}) (string, error) {
	thisString := ""
	switch param.(type) {
	case string:
		if value, ok := param.(string); ok {
			thisString = value
		}
	case float64:
		if value, ok := param.(float64); ok {
			thisString = strconv.FormatFloat(value, 'f', -1, 64)
		}
	case int:
		if value, ok := param.(int); ok {
			thisString = strconv.Itoa(value)
		}
	case int64:
		if value, ok := param.(int64); ok {
			thisString = strconv.FormatInt(value, 10)
		}
	case float32:
		if value, ok := param.(float32); ok {
			thisString = strconv.FormatFloat(float64(value), 'f', -1, 32)
		}
	case bool:
		if value, ok := param.(bool); ok {
			if value == true {
				thisString = "true"
			} else {
				thisString = "false"
			}
		}
	default:
		return "", errors.New("unknow type")
	}

	return thisString, nil
}
