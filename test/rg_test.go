package rgtest

import (
	"github.com/jackylee92/rgo"
	"github.com/jackylee92/rgo/util/rgstring"
	"log"
	"testing"
	"unsafe"
)

// go test -v -run TestStringByte test/rg_test.go -count=1 -config=../config.yaml
func TestStringByte(t *testing.T) {
	str := "abcdecg"
	aa := *(*uintptr)(unsafe.Pointer(uintptr(unsafe.Pointer(&str)) + uintptr(0)))
	log.Println("==", &str, aa)
	b := rgstring.StrToByte(&str)
	aa = *(*uintptr)(unsafe.Pointer(uintptr(unsafe.Pointer(b)) + uintptr(0)))
	log.Println("==", b, aa)
	newStr := rgstring.ByteToStr(b)
	aa = *(*uintptr)(unsafe.Pointer(uintptr(unsafe.Pointer(newStr)) + uintptr(0)))
	log.Println("==", newStr, aa)
	rgo.This.Log.Info("结果", str, b, newStr)
}

// go test -v -run TestStart test/rg_test.go -count=1 -config=../config.yaml
//func TestStart(t *testing.T) {
//	rgo.This.Log.Info("test", "start")
//}
