package rgenv

import "github.com/jackylee92/rgo/core/rgconfig"

func GetEnv() string {
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