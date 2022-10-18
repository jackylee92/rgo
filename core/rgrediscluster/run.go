package rgrediscluster

import (
	"errors"
	"github.com/go-redis/redis"
	"github.com/jackylee92/rgo/core/rglog"
	"strings"
	"sync"
)

// <LiJunDong : 2022-03-30 21:21:22> --- redis配置
type Config struct {
	Addrs        []string // 集群链接地址
	Password     string   // 链接认证
	PoolSize     int      // 链接池大小
	MinIdleConns int      // 初始化连接数
	DB           int      // DB
	TimeOut      int64
}

var redisClientPool sync.Map

type Client struct {
	linkObj *redis.ClusterClient
}

// Start 启动 集群模式，返回key，通过key获取链接对象，调用redis方法执行
// @Param   :
// @Return  :
// @Author  : LiJunDong
// @Time    : 2022-05-31
func Start(cfg Config) (key string, err error) {
	if err = validateConfig(cfg); err != nil {
		return key, err
	}
	// <LiJunDong : 2022-03-30 21:52:55> --- 集群
	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:        cfg.Addrs,
		Password:     cfg.Password,
		PoolSize:     cfg.PoolSize,     // 连接池最大socket连接数，默认为10倍CPU数， 10 * runtime.NumCPU
		MinIdleConns: cfg.MinIdleConns, //在启动阶段创建指定数量的Idle连接，并长期维持idle状态的连接数不少于指定数量；。
	})
	err = client.Ping().Err()
	if err != nil {
		panic("redis链接失败｜" + err.Error())
	}
	if err = client.Ping().Err(); err != nil {
		panic("redis链接失败|" + strings.Join(cfg.Addrs, ","))
	}
	rglog.SystemInfo("启动项【redis】", strings.Join(cfg.Addrs, ","), "成功")
	key = strings.ReplaceAll(strings.ReplaceAll(strings.Join(cfg.Addrs, ""), ".", ""), ":", "")
	redisClientPool.Store(key, &Client{
		linkObj: client,
	})
	return key, err
}

func Pool(key string) (client *Client, err error) {
	clientITF, ok := redisClientPool.Load(key)
	if !ok {
		return client, errors.New("获取链接对象失败")
	}
	client = clientITF.(*Client)
	return client, err
}

/*
* @Content : 获取根据不同配置文件方式redis链接配置
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2022-05-12
 */
func validateConfig(cfg Config) (err error) {
	if len(cfg.Addrs) == 0 {
		return errors.New("redis Addrs不能为空")
	}
	return err
}

func (client *Client) GetRedisClient() (baseClient redis.Cmdable, err error) {
	return client.linkObj, nil
}
