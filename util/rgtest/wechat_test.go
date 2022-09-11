package rglibtest

import (
	_ "github.com/jackylee92/rgo"
	"github.com/jackylee92/rgo/core/rgstarthook"
	_ "github.com/jackylee92/rgo/util/rgstarthook"
	"log"
	"testing"
	"time"
)

//go test -v -run Test_startHook util/rgtest/wechat_test.go -count=1 -config=../../config.yaml
func Test_startHook(t *testing.T) {
	log.Println("start--------")
	rgstarthook.Run()
	time.Sleep(3 * time.Second)
}
