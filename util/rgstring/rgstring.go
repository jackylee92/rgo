package rgstring

import (
	"reflect"
	"unsafe"
)

func StrToByte(param *string) (data *[]byte) {
	paramPoint := *(*reflect.StringHeader)(unsafe.Pointer(param))
	data = (*[]byte)(unsafe.Pointer(&paramPoint))
	return data
}

func ByteToStr(param *[]byte) (data *string) {
	paramPoint := *(*reflect.SliceHeader)(unsafe.Pointer(param))
	data = (*string)(unsafe.Pointer(&paramPoint))
	return data
}
