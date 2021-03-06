package rgcache

import (
	"github.com/coocood/freecache"
	"github.com/go-redis/redis"
	"rgo/core/rgconfig"
	"rgo/core/rgglobal/rgconst"
	"rgo/core/rglog"
	"rgo/core/rgmodel/rgredis"
)

/*
 * @Content : rgredis 同步数据
 * @Author  : LiJunDong
 * @Time    : 2022-05-31$
 */

const (
	memorySize = "util_cache_memory_size" // freecache 申请内存大小
	redisKeyPrefix = "util_cache_listen_redis_key_prefix" // 监听的redis中key的前缀
)
var Client *freecache.Cache // <LiJunDong : 2022-06-01 14:53:33> --- 对外开放，项目中可以使用本地缓存

func init() {
	run()
}
// run 启动
// @Param   :
// @Return  :
// @Author  : LiJunDong
// @Time    : 2022-05-30
func run()  {
	clientType := rgredis.GetClientType()
	// <LiJunDong : 2022-06-01 11:26:17> --- 1: 单机 2：集群
	if clientType != 1 && clientType != 2 {
		rglog.SystemError("redis链接错误")
		return
	}
	clientInterface, err := rgredis.GetClient()
	if err != nil {
		rglog.SystemError("开启监听失败，获取链接失败|" + err.Error())
		return
	}
	configSize := rgconfig.GetInt(memorySize)
	if configSize == 0 {
		rglog.SystemError("freecahce启动失败，未设置容量|" + memorySize)
		return
	}
	freecacheUp(int(configSize))
	go listen(clientType, clientInterface)
}

// listen 订阅redis
// @Param   :
// @Return  :
// @Author  : LiJunDong
// @Time    : 2022-06-01
func listen(clientType int, clientInterface redis.Cmdable) () {
	keyPrefix := rgconfig.GetStr(redisKeyPrefix)
	redisDB := rgconfig.GetStr(rgconst.ConfigKeyRedisDB)
	if redisDB == "" {
		redisDB = "0"
	}
	baseChanName := "__keyspace@" +redisDB+"__:" // <LiJunDong : 2022-06-01 17:41:18> --- 取配置文件中sys_redis_db, 默认为0
	channelName := baseChanName
	if keyPrefix != "" {
		channelName += keyPrefix+ "*"
	}
	var pubsub *redis.PubSub
	if clientType == 1 {
		client := clientInterface.(*redis.Client)
		pubsub = client.PSubscribe(channelName)
	}
	if clientType == 2 {
		client := clientInterface.(*redis.ClusterClient)
		pubsub = client.PSubscribe(channelName)
	}
	ch := pubsub.Channel()
	for msg := range ch {
		// <LiJunDong : 2022-06-01 13:48:12> --- msg.String(): Message<__keyspace@0__:rgo:demo_cron:1: set>
		// <LiJunDong : 2022-06-01 13:48:38> --- msg.Channel: __keyspace@0__:rgo:demo_cron:1
		// <LiJunDong : 2022-06-01 13:48:38> --- msg.Payload: set
		// <LiJunDong : 2022-06-01 13:48:38> --- msg.Pattern: msg.Payload
		//log.Println("listen message", msg.String(), msg.Channel, msg.Payload, msg.Pattern)
		key := msg.Channel[len(baseChanName):]
		action := msg.Payload
		del(action, key)
	}
}

// del
// @Param   : key string
// @Return  :
// @Author  : LiJunDong
// @Time    : 2022-06-01
func del(action string, keys ...string) () {
	for _, item := range keys {
		Client.Del([]byte(item))
	}
}


// freecacheUp 启动freecahche
// @Param   :
// @Return  :
// @Author  : LiJunDong
// @Time    : 2022-06-01
func freecacheUp(size int)  {
	cacheSize := size * 1024 * 1024
	Client = freecache.NewCache(cacheSize)
}