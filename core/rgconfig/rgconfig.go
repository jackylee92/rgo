package rgconfig

import (
	"github.com/jackylee92/rgo/core/rgglobal/rgerror"
	"io/ioutil"
	"os"
	"strings"

	"flag"

	"gopkg.in/yaml.v2"
)

/* 说明 <LiJunDong : 2022-03-02 14:05:42> --- 初始化获取配置
 *
 */

// <LiJunDong : 2022-03-10 14:33:29> --- 配置集合
var configSource = make(map[string]ConfigInterface)

// 目前使用的配置
var configPool ConfigInterface

/*
* @Content : register
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2022-03-11
 */
func Register(name string, c ConfigInterface) {
	configSource[name] = c
}

/*
* @Content : 设置
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2022-03-09
 */
func getFactoryClient(param string) ConfigInterface {
	return configSource[param]
}

/*
* @Content : 加载配置
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2022-03-01
 */
func Start() {
	loadBaseConfig()
	configPool = getFactoryClient(Config.Config)
	if configPool == nil {
		panic("配置启动失败，无[" + Config.Config + "]方式配置")
	}
	if configPool.Load() == nil {
		panic("配置启动失败，无[" + Config.Config + "]加载方法为空")
	}
	configPool.Load()()
}

type BaseConfig struct {
	Config               string `yaml:"config"`
	ApolloAppId          string `yaml:"apollo_app_id"`
	ApolloCluster        string `yaml:"apollo_cluster"`
	ApolloHost           string `yaml:"apollo_host"`
	ApolloIsBackupConfig bool   `yaml:"apollo_is_backup_config"`
	ApolloSecret         string `yaml:"apollo_secret"`
	ApolloNamespaceName  string `yaml:"apollo_namespace_name"`
	SysAppName           string `yaml:"sys_app_name"`
}
type Option struct {
	Config string `long:"config" required:"false"`
}

var (
	configPath = flag.String("config", "", "基础配置")
)

var configByte []byte
var Config BaseConfig

/*
* @Content : 获取配置文件字节
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2022-05-04
 */
func loadBaseConfig() {
	// <LiJunDong : 2022-05-12 22:46:19> --- 适配args为config需要的格式
	oldArgs := parse()
	flag.Parse()
	// <LiJunDong : 2022-05-12 22:46:19> --- 还原原来的args
	os.Args = oldArgs
	if *configPath == "" {
		panic(rgerror.ErrorStartConfigExists)
	}
	if _, err := os.Stat(*configPath); os.IsNotExist(err) {
		panic(rgerror.ErrorConfigExist + ":" + *configPath)
	}
	var err error
	configByte, err = ioutil.ReadFile(*configPath)
	if err != nil {
		panic(rgerror.ErrorBaseConfig + ":" + err.Error())
	}
	err = yaml.Unmarshal(configByte, &Config)
	if err != nil {
		panic(rgerror.ErrorConfigParse + ":" + err.Error())
	}
	if Config.Config != "apollo" && Config.Config != "file" {
		panic( rgerror.ErrorConfigTypeExist + ":" + Config.Config)
	}
	return
}

/*
* @Content : 获取配置byte字节
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2022-05-11
 */
func ReadFile() []byte {
	return configByte
}

/*
* @Content : 处理命令行参数
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2022-05-12
 */
func parse() []string {
	oldArgs := os.Args
	newArgs := make([]string, 0, 2)
	if len(os.Args) > 0 {
		newArgs = append(newArgs, os.Args[0])
	}
	for key, item := range os.Args {
		if key == 0 {
			continue
		}
		if strings.Contains(item, "-config") {
			newArgs = append(newArgs, item)
			break
		}
	}
	if len(newArgs) == 2 {
		os.Args = newArgs
	}
	return oldArgs
}
