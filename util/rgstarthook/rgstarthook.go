package rgstarthook

import (
	"github.com/jackylee92/rgo/core/rgconfig"
	"github.com/jackylee92/rgo/core/rgglobal"
	"github.com/jackylee92/rgo/core/rgstarthook"
	"github.com/jackylee92/rgo/util/rgmsg"
)

const (
	configNoticeWechat = "util_start_msg_notice_wechat"
	configNoticeEmailTo = "util_start_msg_notice_email_to"
	configNoticeEmailFrom = "util_start_msg_notice_email_from"
)

func init() {
	rgstarthook.RegisterStartFunc(startNotice)
}

/*
 * @Content : rgstart
 * @Author  : LiJunDong
 * @Time    : 2022-05-30$
 */

// startNotice 启动通知
// @Param   :
// @Return  :
// @Author  : LiJunDong
// @Time    : 2022-05-30
func startNotice()  {
    rgmsgClient := rgmsg.Client{
    	This: nil,
    	Title: "服务启动通知",
	}
	wechatOk := rgconfig.GetBool(configNoticeWechat)
	if wechatOk {
		rgmsgClient.Typ = rgmsg.WECHAT
		rgmsgClient.WechatContent = "服务启动通知"
		rgmsgClient.WechatData = []string{"server: " + rgglobal.AppName}
		rgmsgClient.Send()
	}
	toEmailList := rgconfig.GetStrSlice(configNoticeEmailTo)
	fromEmailList := rgconfig.GetStr(configNoticeEmailFrom)
	if len(toEmailList) != 0 && fromEmailList != "" {
		rgmsgClient.Typ = rgmsg.EMAIL
		rgmsgClient.EmailTo = toEmailList
		rgmsgClient.EmailContent = "服务启动通知【" +rgglobal.AppName+ "】"
		rgmsgClient.EmailFrom = fromEmailList
		rgmsgClient.Send()
	}
	return
}

// getEnvValue 获取环境
// @Param   :
// @Return  : string
// @Author  : LiJunDong
// @Time    : 2022-05-30
func getEnvValue() string {
	env := rgconfig.GetStr("env")
	value := "未知环境"
	switch env {
	case "local":
		value = "本地环境"
	case "dev":
		value = "开发环境"
	case "test":
		value = "测试环境"
	case "uat":
		value = "uat环境"
	case "prod":
		value = "生产环境"
	}
	return value
}