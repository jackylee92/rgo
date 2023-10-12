package rgemail

import (
	"errors"
	"github.com/jackylee92/rgo/core/rgrequest"
	"gopkg.in/gomail.v2"
)

type EmailClient struct {
	This        *rgrequest.Client `json:"-"` // json解析忽略
	Title       string            //标题
	To          []string          //接收方
	Cc          []string          //抄送方
	BCc         []string          //秘密抄送方
	Content     string            //邮件内容
	Attach      string            //邮件附件路径
	From        string            //发送方 默认：system@ruigushop.com
	ContentType string            //邮件内容格式 默认text/html
	Host        string            //邮件服务器的地址 默认使用system@ruigushop.com
	Port        int               //邮件服务器端口 465
	UserName    string            //用户名
	Password    string            //密码
}

func (e *EmailClient) getEmailBodyType() string {
	return "text/html"
}
func (e *EmailClient) getHost() string {
	if e.Host == "" {
		return "smtp.mxhichina.com"
	}
	return e.Host
}
func (e *EmailClient) getPort() int {
	if e.Port <= 0 {
		return 465
	}
	return e.Port
}
func (e *EmailClient) getUserName() string {
	if e.UserName == "" {
		return "system@ruigushop.com"
	}
	return e.UserName
}

func (e *EmailClient) getFrom() string {
	if e.From == "" {
		return "system@ruigushop.com"
	}
	return e.From
}

func (e *EmailClient) getPassword() string {
	if e.Password == "" {
		return "oxedYa3sQckk2f"
	}
	return e.Password
}

func (e *EmailClient) Send() error {
	var err error
	err = e.check()
	if err != nil {
		return err
	}
	e.ContentType = e.getEmailBodyType()
	m := gomail.NewMessage()
	m.SetHeader("From", e.getFrom())
	m.SetHeader("To", e.To...) //主送
	if len(e.Cc) != 0 {
		m.SetHeader("Cc", e.Cc...) //抄送
	}
	if len(e.BCc) != 0 {
		m.SetHeader("Bcc", e.BCc...) // 密送
	}
	m.SetHeader("Subject", e.Title)
	m.SetBody(e.ContentType, e.Content)
	if e.Attach != "" {
		m.Attach(e.Attach) //添加附件
	}
	d := gomail.NewDialer(e.getHost(), e.getPort(), e.getUserName(), e.getPassword())
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

func (e *EmailClient) check() error {

	if e.Content == "" {
		return errors.New("请配置邮件内容")
	}
	if e.Title == "" {
		return errors.New("请配置邮件标题")
	}
	if e.Host == "" {
		return errors.New("请配置邮件Host")
	}
	if len(e.To) <= 0 {
		return errors.New("请配置发送者")
	}

	return nil
}
