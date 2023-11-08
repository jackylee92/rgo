package rgenv

import (
	"github.com/jackylee92/rgo/core/rgconfig"
	"github.com/jackylee92/rgo/core/rgglobal/rgconst"
	"net"
	"os"
)

var hostName string
var appName string
var env string
var ipStr string

func GetEnv() string {
	if env == "" {
		env = rgconfig.GetStr("env")
	}
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
func GetMachine() string {
	if hostName == "" {
		hostName, _ = os.Hostname()
	}
	return hostName
}

func GetIp() string {
	if ipStr == "" {
		addrs, err := net.InterfaceAddrs()
		if err == nil {
			for _, address := range addrs {
				if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						ipStr = ipnet.IP.String()
					}

				}
			}
		}
	}
	return ipStr
}

func GetAppName() string {
	if appName == "" {
		appName = rgconfig.GetStr(rgconst.ConfigKeyAppName)
	}
	return appName
}
