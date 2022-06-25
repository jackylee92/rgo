package rgjson

import (
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func Marshel(data interface{}) (result string, err error) {
	resultByte, err := json.Marshal(data)
	if err != nil {
		return result, err
	}
	result = string(resultByte)
	return result, err
}

func UnMarshel(data []byte, v interface{}) (err error) {
	err = json.Unmarshal(data, v)
	if err != nil {
		return err
	}
	return err
}
