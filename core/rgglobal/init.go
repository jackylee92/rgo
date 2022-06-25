package rgglobal

import (
	"net"
	"os"

	"rgo/core/rgglobal/rgerror"
)

var BasePath string

// var BasePathConfig string
var AppName string
var LocalIp string

/*
* @Content :
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2022-03-02
 */
func Load() {
	// <LiJunDong : 2022-03-08 18:28:56> --- 设置项目所在绝对目录
	setBasePath()
	// BasePathConfig = BasePath + rgconst.BaseConfigEnd
	setLocalIP()
}

/*
* @Content : 设置程序根目录
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2022-03-08
 */
func setBasePath() {
	basePath, err := os.Getwd()
	if err != nil {
		panic(rgerror.ErrorBasePath + "|err" + err.Error())
	}
	BasePath = basePath
}

func SetAppName(name string) {
	if name == "" {
		panic(rgerror.ErrorAppNameNil)
	}
	AppName = name
}

func setLocalIP() {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		return
	}
	for i := 0; i < len(netInterfaces); i++ {
		if (netInterfaces[i].Flags & net.FlagUp) != 0 {
			addrs, _ := netInterfaces[i].Addrs()
			for _, address := range addrs {
				if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					//获取IPv6
					/*if ipnet.IP.To16() != nil {
					    ipStr = append(ipStr, ipnet.IP.String())

					}*/
					//获取IPv4
					if ipnet.IP.To4() != nil {
						LocalIp = ipnet.IP.String()
						return
					}
				}
			}
		}
	}
	return
}
