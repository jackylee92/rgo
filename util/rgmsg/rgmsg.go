package rgmsg

import (
	"errors"
	"github/tidwall/gjson"
	"gopkg.in/gomail.v2"
	_ "rgo"
	"rgo/core/rgconfig"
	"rgo/core/rgjson"
	"rgo/core/rgrequest"
	"rgo/util/rghttp"
	"rgo/util/rgtime"
)

/*
 * @Content : <LiJunDong : 2022-05-27 12:05:25> --- GitHub: https://github/go-gomail/gomail
 * @Author  : LiJunDong
 * @Time    : 2022-05-27$
 */

type msgType int

const (
	EMAIL msgType = iota + 1
	WECHAT
)

const (
	configEmailHost  = "util_email_host"  // <LiJunDong : 2022-05-28 15:26:41> --- 发送邮件服务器地址
	configEmailPort  = "util_email_port"  // <LiJunDong : 2022-05-28 15:26:41> --- 发送邮件服务端口
	configEmailLogin = "util_email_login" // <LiJunDong : 2022-05-28 15:27:28> --- 发送邮件用户名
	configEmailPwd   = "util_email_pwd"   // <LiJunDong : 2022-05-28 15:27:28> --- 发送邮件用户密码
	configWechatTo   = "util_wechat_to"   // <LiJunDong : 2022-05-28 15:20:08> --- 企业微信接收方
)

type Client struct {
	This          *rgrequest.Client `json:"-"` // json解析忽略
	Typ           msgType
	Title         string // <LiJunDong : 2022-05-27 12:00:36> --- 公用标题
	EmailTo       []string // <LiJunDong : 2022-05-27 12:01:44> --- 接受方
	EmailCc       []string // <LiJunDong : 2022-05-27 12:01:39> --- 抄送方
	EmailBCc      string // <LiJunDong : 2022-05-27 12:01:32> --- 秘密抄送方
	EmailContent  string // <LiJunDong : 2022-05-27 12:02:06> --- 邮件内容
	EmailAttach   string // <LiJunDong : 2022-05-27 12:01:32> --- 邮件附件路径
	EmailFrom     string // <LiJunDong : 2022-05-28 15:00:08> --- 发送方
	emailBodyType string // <LiJunDong : 2022-05-27 12:02:41> --- 邮件内容格式 默认text/html

	WechatContent string   // <LiJunDong : 2022-05-27 12:08:21> --- 企业微信消息中 主要内容描述
	WechatData    []string // <LiJunDong : 2022-05-27 12:08:51> --- 企业微信中key:value形式展示在内容描述下面
}

var (
	chanMax    = make(chan struct{}, 10000)
	emailHost  string
	emailPort  int
	emailLogin string
	emailPwd   string
	wechatTo   string
)

func init() {
	setConfig()
}

// Send 发送邮件
// @Param   :
// @Return  : err error 一般情况，异常邮件不需要考虑返回值，代码中可以_忽略
// @Author  : LiJunDong
// @Time    : 2022-05-27
func (c *Client) Send() (err error) {
	if err = c.check(); err != nil {
		return err
	}
	chanMax <- struct{}{}
	defer func() {
		<-chanMax
	}()
	if c.Typ == EMAIL {
		err = c.sendEmail()
	}
	if c.Typ == WECHAT {
		err = c.sendWechat()
	}
	return err
}

// check 检查发送参数
// @Param   :
// @Return  : err error
// @Author  : LiJunDong
// @Time    : 2022-05-27
func (c *Client) check() (err error) {
	if c.Typ == 0 {
		return errors.New("发送方式Typ不能为空")
	}
	if c.Typ != EMAIL && c.Typ != WECHAT {
		return errors.New("发送方式只支持EMAIL/WECHAT")
	}
	if c.Title == "" {
		return errors.New("消息标题不能为空")
	}
	titleRune := ([]rune)(c.Title)
	if len(titleRune) > 100 {
		return errors.New("标题不能超过80个字符")
	}
	if c.Typ == EMAIL {
		if c.EmailFrom == "" {
			return errors.New("邮件配置发送方[EmailFrom]为空")
		}
		if len(c.EmailTo) == 0 {
			return errors.New("邮件配置接受方[EmailTo]为空")
		}
		if c.EmailContent == "" {
			return errors.New("邮件内容[EmailContent]为空")
		}
		if emailHost == "" || emailPort == 0 || emailLogin == "" || emailPwd == "" {
			return errors.New("邮件服务配置错误")
		}
	}
	if c.Typ == WECHAT {
		if wechatTo == "" {
			return errors.New("企业微信地址[" +configWechatTo+ "]配置错误")
		}
		if c.WechatContent == "" {
			return errors.New("企业微信消息内容[WechatContent]不能为空")
		}
	}
	return err
}

// sendWechat 发送企业微信信息
// @Param   :
// @Return  : err error
// @Author  : LiJunDong
// @Time    : 2022-05-27
func (c *Client) sendWechat() (err error) {
	content := "# <font color=\"warning\">" + c.Title + "</font>\n"
	content += "<font color=\"comment\" size=\"10\">at " + rgtime.NowDateTime() + "</font>\n"
	content += c.WechatContent + "\n"
	if len(c.WechatData) != 0 {
		for _, value := range c.WechatData{
			content += "> <font color=\"comment\" size=\"11\">" + value + "</font> \n"
		}
	}
	param := map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]string{
			"content": content,
		},
	}
	paramJson, _ := rgjson.Marshel(param)
	httpClient := rghttp.Client{
		Url:    wechatTo,
		Method: "POST",
		Header: c.getWechatHeader(),
		Param:  paramJson,
		This: c.This,
	}
	data, err := httpClient.GetApi()
	if err != nil {
		return err
	}
	if gjson.Get(data, "code").Int() != 0 {
		return errors.New(gjson.Get(data, "errmsg").String())
	}
	return err
}

// sendEmail 发送邮件
// @Param   :
// @Return  : err error
// @Author  : LiJunDong
// @Time    : 2022-05-27
func (c *Client) sendEmail() (err error) {
	c.emailBodyType = c.getEmailBodyType()
	m := gomail.NewMessage()
	m.SetHeader("From", c.EmailFrom)
	m.SetHeader("To", c.EmailTo...) //主送
	if len(c.EmailCc) != 0 {
		m.SetHeader("Cc", c.EmailCc...) //抄送
	}
	if c.EmailBCc != "" {
		m.SetHeader("Bcc", c.EmailBCc) // 密送
	}
	m.SetHeader("Subject", c.Title)
	m.SetBody(c.emailBodyType, c.EmailContent)
	if c.EmailAttach != "" {
		m.Attach(c.EmailAttach) //添加附件
	}
	d := gomail.NewDialer(emailHost, emailPort, emailLogin, emailPwd)
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return err
}

// setConfig 设置配置，从配置文件添加到全局变量中
// @Param   :
// @Return  :
// @Author  : LiJunDong
// @Time    : 2022-05-28
func setConfig() {
	emailHost = rgconfig.GetStr(configEmailHost)
	emailPort = int(rgconfig.GetInt(configEmailPort))
	emailLogin = rgconfig.GetStr(configEmailLogin)
	emailPwd = rgconfig.GetStr(configEmailPwd)

	wechatTo = rgconfig.GetStr(configWechatTo)
}

// setEmailBodyType 设置邮件内容格式
// @Param   :
// @Return  :
// @Author  : LiJunDong
// @Time    : 2022-05-28
func (c *Client) getEmailBodyType() string {
	return "text/html"
}

// getWechatHeader 获取企业微信消息头部
// @Param   :
// @Return  : map[string]string
// @Author  : LiJunDong
// @Time    : 2022-05-28
func (c *Client)getWechatHeader() map[string]string {
	return map[string]string{
		"Content-Type": "application/json",
	}
}