package rglibtest

import (
	"fmt"
	"github.com/jackylee92/rgo"
	"github.com/jackylee92/rgo/library/rguser"
	"github.com/magiconair/properties/assert"
	"testing"
)

//go test -v -run library/rglibtest user_test.go -count=1 -args -config=../../config.yaml
func Test_checkGetUserInfoByToken(t *testing.T) {

	type args struct {
		UserId int
		Token  string
	}

	tests := []args{
		//地推b
		{UserId: 4860, Token: "8ae5a2c72f7280c50d3857aaa542ffba"},
		{UserId: 61813, Token: "4c90481440e8c0ae1403053b0af8b4b3"},
		//全国b
		{UserId: 4120, Token: "f14dbf1d24ad13097d4a06cbb2883a9d"},
		//RDC
		{UserId: 320, Token: "1f3151592af76b28605c010c2a2caa0d"},
		//C端
		{UserId: 471867, Token: "7ef43f69a159e9425ff2750d34f449fc"},
	}
	for _, tt := range tests {
		userInfo, err := rguser.GetUserInfoBaseByToken(rgo.This, tt.UserId, tt.Token)
		if err != nil {
			t.Errorf("GetUserInfoBaseByToken 返回异常 %v", err)
		}
		if userInfo.Id == 0 {
			t.Errorf("GetUserInfoBaseByToken 未获取到信息")
		}
		assert.Equal(t, len(userInfo.Mobile) > 0, true, "手机号不存在")
	}
}

//go test -v -run library/rglibtest user_test.go -count=1 -args -config=../../config.yaml
func Test_checkGetUserInfoById(t *testing.T) {
	tests := [5]int{4860, 61813, 4120, 320, 471867}

	for _, tt := range tests {
		userInfo, err := rguser.GetUserInfoById(rgo.This, tt)
		assert.Equal(t, err == nil, true, fmt.Sprintf("获取用户信息异常: %v", err))
		assert.Equal(t, userInfo.Id == 0, false, "未获取到用户信息")
		assert.Equal(t, len(userInfo.Mobile) > 0, true, "手机号不存在")
	}

}
