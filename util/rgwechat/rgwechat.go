package rgwechat

import (
	"errors"
	"fmt"
	"github.com/jackylee92/rgo/core/rgjson"
	"github.com/jackylee92/rgo/core/rgrequest"
	"github.com/jackylee92/rgo/util/rgenv"
	"github.com/jackylee92/rgo/util/rghttp"
	"github.com/jackylee92/rgo/util/rgtime"
	"github.com/tidwall/gjson"
)

type msgType int

const (
	Error msgType = iota + 1
	Info
)

type WeClient struct {
	This    *rgrequest.Client `json:"-"` // json解析忽略
	Content []string
	Title   string
	To      string
	Level   msgType
}

func (c *WeClient) Send() (bool, error) {
	err := c.checkConfig()
	if err != nil {
		return false, err
	}
	param := map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]string{
			"content": c.getContent(),
		},
	}
	paramJson, _ := rgjson.Marshel(param)
	httpClient := rghttp.Client{
		Url:    c.To,
		Method: "POST",
		Header: c.getHeader(),
		Param:  paramJson,
		This:   c.This,
	}
	data, err := httpClient.GetApi()
	if err != nil {
		fmt.Println("err2", err)
		return false, err
	}
	if gjson.Get(data, "errcode").Int() != 0 {
		return false, errors.New("发送企业微信消息失败：" + gjson.Get(data, "errmsg").String())
	}
	return true, nil
}

func (c *WeClient) getContent() (content string) {
	msgType := "warning"
	if c.Level == Info {
		msgType = "info"
	}
	content = "# <font color=\"comment\" size=\"11\">【 " + rgenv.GetEnv() + "】</font>\n"
	content += "# <font size=\"16\" color=\"" + msgType + "\"> " + c.Title + "</font>\n"
	content += "<font color=\"comment\" size=\"10\">at " + rgtime.NowTime() + "</font>\n"
	if len(c.Content) != 0 {
		for _, value := range c.Content {
			content += "> <font color=\"comment\" size=\"13\">" + value + "</font> \n"
		}
	}
	return content
}

func (c *WeClient) checkConfig() error {
	if c.To == "" {
		return errors.New("未配置机器人")
	}
	if c.Title == "" {
		return errors.New("请传入标题")
	}
	if len(c.Content) <= 0 {
		return errors.New("发送的内容不能为空")
	}
	return nil
}

func (c *WeClient) getHeader() map[string]string {
	return map[string]string{
		"Content-Type": "application/json",
	}
}
