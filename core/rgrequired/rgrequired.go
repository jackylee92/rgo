package rgrequired

import (
	"github.com/jackylee92/rgo/core/rgconfig"
	"github.com/jackylee92/rgo/core/rgglobal/rgconst"
	"github.com/jackylee92/rgo/core/rglog"
)

// <LiJunDong : 2022-03-02 14:10:55> --- 验证框架必须的目录结构

/*
* @Content : 验证目录
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2022-03-02
 */
func checkRequiredFolders() {
	// if _, err := os.Stat(rgglobal.BasePathConfig); err != nil {
	//     rglog.SystemError(error_enum.ErrorConfigExist + err.Error())
	// }
}

func Check() {
	// checkRequiredFolders()
	checkConfig()
}

/*
* @Content : 检查配置文件是否齐全
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2022-03-28
 */
func checkConfig() {
	value := rgconfig.GetStr(rgconst.ConfigKeyLogLevel)
	// log.Println(rgconst.ConfigKeyLogLevel, value)
	// if value == "" {
	//     rglog.SystemInfo("[" + rgconst.ConfigKeyLogLevel + "]为空")
	// }
	value = rgconfig.GetStr(rgconst.ConfigKeyMysql)
	// log.Println(rgconst.ConfigKeyMysql, value)
	// if value == "" {
	//     rglog.SystemInfo("[" + rgconst.ConfigKeyMysql + "]为空")
	// }
	value = rgconfig.GetStr(rgconst.ConfigKeyAppName)
	// log.Println(rgconst.ConfigKeyAppName, value)
	if value == "" {
		// rglog.SystemInfo("[" + rgconst.ConfigKeyAppName + "]为空")
		panic("[" + rgconst.ConfigKeyAppName + "]为空")
	}
	value = rgconfig.GetStr(rgconst.ConfigKeyPort)
	// log.Println(rgconst.ConfigKeyPort, value)
	if value == "" {
		// rglog.SystemInfo("[" + rgconst.ConfigKeyPort + "]为空")
		panic("[" + rgconst.ConfigKeyPort + "]为空")
	}
	value = rgconfig.GetStr(rgconst.ConfigKeyMessage)
	// log.Println(rgconst.ConfigKeyMessage, value)
	if value == "" {
		rglog.SystemInfo("[" + rgconst.ConfigKeyMessage + "]为空")
	}
	value = rgconfig.GetStr(rgconst.ConfigKeyJaegerStatus)
	// log.Println(rgconst.ConfigKeyJaegerStatus, value)
	value = rgconfig.GetStr(rgconst.ConfigKeyJaegerHost)
	// log.Println(rgconst.ConfigKeyJaegerHost, value)
	// if value == "" {
	//     rglog.SystemInfo("[" + rgconst.ConfigKeyJaegerHost + "]为空")
	// }
	value = rgconfig.GetStr(rgconst.ConfigKeyRequestLog)
	// log.Println(rgconst.ConfigKeyRequestLog, value)
	if value == "" {
		rglog.SystemInfo("[" + rgconst.ConfigKeyRequestLog + "]为空")
	}
}
