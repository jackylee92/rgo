package rgapollo

import (
	"github.com/jackylee92/rgo/core/rgconfig"

	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/env/config"
	"github.com/apolloconfig/agollo/v4/storage"
)

const (
	Name string = "apollo"
)

type ApolloConfig struct {
	AppId          string `json:"app_id"`
	Cluster        string `json:"cluster"`
	Host           string `json:"host"`
	IsBackupConfig bool   `json:"is_backup_config"`
	Secret         string `json:"secret"`
	NamespaceName  string `json:"namespace_name"`
}

var configPool *storage.Config

/* 从apollo获取：先从appollo获取，获取失败则从本地备份配置文件获取，获取到更新本地备份配置文件
 *
 */

func Register() {
	rgconfig.Register(Name, client{})
}

func (c client) getConfig() (config ApolloConfig) {
	config = ApolloConfig{
		AppId:          rgconfig.Config.ApolloAppId,
		Cluster:        rgconfig.Config.ApolloCluster,
		Host:           rgconfig.Config.ApolloHost,
		IsBackupConfig: rgconfig.Config.ApolloIsBackupConfig,
		Secret:         rgconfig.Config.ApolloSecret,
		NamespaceName:  rgconfig.Config.ApolloNamespaceName,
	}

	return config
}

func (c client) load() (err error) {
	apolloConfig := c.getConfig()
	apolloC := &config.AppConfig{
		AppID:          apolloConfig.AppId,
		Cluster:        apolloConfig.Cluster,
		IP:             apolloConfig.Host,
		IsBackupConfig: apolloConfig.IsBackupConfig,
		NamespaceName:  apolloConfig.NamespaceName,
		Secret:         apolloConfig.Secret,
	}
	client, err := agollo.StartWithConfig(func() (*config.AppConfig, error) {
		return apolloC, nil
	})
	if err != nil {
		panic("初始化Apollo配置失败|" + err.Error())
	}

	configPool = client.GetConfig(apolloC.NamespaceName)
	return err
}
