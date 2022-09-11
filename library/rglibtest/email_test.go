package rglibtest

import (
	"fmt"
	"github.com/jackylee92/rgo"
	"github.com/jackylee92/rgo/library/rgemail"
	"github.com/magiconair/properties/assert"
	"testing"
)

//go test -v -run library/rglibtest email_test.go -count=1 -args -config=../../config.yaml
func Test_sendEmail(t *testing.T) {

	type emailParam struct {
		Cc      []string
		To      []string
		Content string
		Title   string
	}

	tests := []emailParam{
		{Cc: []string{"wangxiaoyang@ruigushop.com"}, To: []string{"wangxiaoyang@ruigushop.com"}, Content: "this is demo", Title: "demo"},
	}

	for _, tt := range tests {
		client := rgemail.EmailClient{
			This: rgo.This,
			Cc:      tt.Cc,
			To:      tt.To,
			Content: tt.Content,
			Title: "demo",
		}
		err := client.Send()
		assert.Equal(t, err == nil, true, fmt.Sprintf("发送邮件异常： %v", err))
	}

}
