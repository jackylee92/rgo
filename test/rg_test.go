package rgtest

import (
	"github.com/jackylee92/rgo"
	"github.com/jackylee92/rgo/core/rglog"
	"github.com/jackylee92/rgo/core/rgrequest"
	"testing"
	"time"
)

// go test -v -run TestLog test/rg_test.go -count=1 -config=../config.yaml
func TestLog(t *testing.T) {
	this := rgo.This
	go wLog(this)
	time.Sleep(180 * time.Second)
	return
}

func wLog(this *rgrequest.Client) {
	for {
		time.Sleep(1 * time.Second)
		this.Log.Debug("d级别日志", "ddddd")
		this.Log.Info("info级别日志")
		this.Log.Warn("w级别日志")
		this.Log.Error("e级别日志")
		rglog.RequestLog(this.UniqId, rglog.LevelRequest, "请求体")
		rglog.RequestLog(this.UniqId, rglog.LevelResponse, "返回体")
		rglog.SystemInfo("系统级别日志")
		rglog.SystemError("系统级别错误日志")
	}
}
