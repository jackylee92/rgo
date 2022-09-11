package rglibtest

import (
	"fmt"
	"github.com/jackylee92/rgo"
	"github.com/jackylee92/rgo/library/rgwechat"
	"github.com/magiconair/properties/assert"
	"testing"
)

//go test -v -run Test_sendWeChat library/rglibtest/wechat_test.go -count=1 -config=../../config.yaml
func Test_sendWeChat(t *testing.T) {

	type WeChatArgs struct {
		Msg   []string
		Title string
		Want  bool
		To    string
	}

	tests := []WeChatArgs{
		{Msg: []string{"订单号：R0ADJFASQWEXAKSJ", "优惠券编号：DSDDSDCWEESECS"}, Title: "系统消息", Want: true, To: "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=3b578feb-3010-489e-82cf-b20823a52c53"},
		{Msg: []string{"订单号", "R0ADJFASQWEXAKSJ", "优惠券编号：DSDDSDCWEESECS"}, Title: "系统消息", Want: false, To: "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=3b578feb-3010-489e-82cf"},
	}

	for _, tt := range tests {
		client := rgwechat.WeClient{
			This:    rgo.This,
			Content: tt.Msg,
			To:      tt.To,
			Title:   tt.Title,
			Level: rgwechat.Info,
		}
		res, err := client.Send()
		assert.Equal(t, err == nil, tt.Want, fmt.Sprintf("发送出现error: %v", err))
		assert.Equal(t, res == true, tt.Want, fmt.Sprintf("发送失败"))
	}

}
