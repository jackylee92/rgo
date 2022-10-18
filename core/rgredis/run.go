package rgredis

import (
	"errors"
	"github.com/go-redis/redis"
	"github.com/jackylee92/rgo/core/rglog"
	"strings"
	"sync"
	"time"
)

// <LiJunDong : 2022-03-30 21:21:22> --- redis配置
type Config struct {
	Addr         string // 单机链接的地址
	Password     string // 链接认证
	PoolSize     int    // 链接池大小
	MinIdleConns int    // 初始化连接数
	DB           int    // DB
	TimeOut      int64
}

var redisClientPool sync.Map

type Client struct {
	linkObj *redis.Client
}

// Start 启动 单机模式，返回key，通过key获取链接对象，调用redis方法执行
// @Param   :
// @Return  :
// @Author  : LiJunDong
// @Time    : 2022-05-31
func Start(cfg Config) (key string, err error) {
	if err = validateConfig(cfg); err != nil {
		return key, err
	}
	// <LiJunDong : 2022-03-30 21:52:42> --- 单机
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password, // no password set
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,     // 连接池最大socket连接数，默认为10倍CPU数， 10 * runtime.NumCPU
		MinIdleConns: cfg.MinIdleConns, //在启动阶段创建指定数量的Idle连接，并长期维持idle状态的连接数不少于指定数量；。
		ReadTimeout:  time.Duration(cfg.TimeOut) * time.Second,
		WriteTimeout: time.Duration(cfg.TimeOut) * time.Second,
	})
	if err = client.Ping().Err(); err != nil {
		panic("redis链接失败|" + cfg.Addr)
	}
	rglog.SystemInfo("启动项【redis】" + cfg.Addr + ":成功")
	key = strings.ReplaceAll(strings.ReplaceAll(cfg.Addr, ".", ""), ":", "")
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
	if cfg.Addr == "" {
		return errors.New("redis Addr不能为空")
	}
	return err
}

func GetRedisClient(client *Client) (baseClient *redis.Client) {
	return client.linkObj
}
