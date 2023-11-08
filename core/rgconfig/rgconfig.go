package rgconfig

import (
	"github.com/jackylee92/rgo/core/rgglobal/rgerror"
	"os"
	"testing"

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

func Register(name string, c ConfigInterface) {
	configSource[name] = c
}

func getFactoryClient(param string) ConfigInterface {
	return configSource[param]
}

func Start() {
	loadBaseConfig()
	configPool = getFactoryClient(Config.Config)
	if configPool == nil {
		panic("配置启动失败，无[" + Config.Config + "]方式配置")
	}
	if configPool.Load() == nil {
		panic("配置启动失败，无[" + Config.Config + "]加载方法为空")
	}
	_ = configPool.Load()()
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

func loadBaseConfig() {
	// <LiJunDong : 2022-05-12 22:46:19> --- 适配args为config需要的格式
	//oldArgs := parse()
	//flag.Parse()
	//// <LiJunDong : 2022-05-12 22:46:19> --- 还原原来的args
	//
	//os.Args = oldArgs
	testing.Init()
	flag.Parse()
	if *configPath == "" {
		panic(rgerror.ErrorStartConfigExists)
	}
	if _, err := os.Stat(*configPath); os.IsNotExist(err) {
		panic(rgerror.ErrorConfigExist + ":" + *configPath)
	}
	var err error
	configByte, err = os.ReadFile(*configPath)
	if err != nil {
		panic(rgerror.ErrorBaseConfig + ":" + err.Error())
	}
	err = yaml.Unmarshal(configByte, &Config)
	if err != nil {
		panic(rgerror.ErrorConfigParse + ":" + err.Error())
	}
	if Config.Config != "apollo" && Config.Config != "file" {
		panic(rgerror.ErrorConfigTypeExist + ":" + Config.Config)
	}
	return
}

func ReadFile() []byte {
	return configByte
}
