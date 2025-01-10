package rgcache

// 试用easyCache代替
//
//import (
//	"github.com/coocood/freecache"
//	"github.com/go-redis/redis"
//	"github.com/jackylee92/rgo"
//	"github.com/jackylee92/rgo/core/rgconfig"
//	"github.com/jackylee92/rgo/core/rgglobal/rgconst"
//	"github.com/jackylee92/rgo/core/rglog"
//	"github.com/jackylee92/rgo/core/rgredis"
//)
//
///*
// * @Content : rgredis 同步数据
// * @Author  : LiJunDong
// * @Time    : 2022-05-31$
// */
//
//const (
//	memorySize     = "util_cache_memory_size"             // freecache 申请内存大小
//	redisKeyPrefix = "util_cache_listen_redis_key_prefix" // 监听的redis中key的前缀
//)
//
//var Client *freecache.Cache // <LiJunDong : 2022-06-01 14:53:33> --- 对外开放，项目中可以使用本地缓存
//
//type redisClientTypeEnum int
//
//const (
//	RedisEnum redisClientTypeEnum = iota + 1
//	RedisClusterEnum
//)
//
//var redisClient rgredis.RedisClientITF
//var redisClientType redisClientTypeEnum
//
//// run 启动
//// @Param   :
//// @Return  :
//// @Author  : LiJunDong
//// @Time    : 2022-05-30
//func Run(redisClient rgredis.RedisClientITF, enum redisClientTypeEnum) {
//	// <LiJunDong : 2022-06-01 11:26:17> --- 1: 单机 2：集群
//	if enum != 1 && enum != 2 {
//		rglog.SystemError("redisClientType错误")
//		return
//	}
//	configSize := rgconfig.GetInt(memorySize)
//	if configSize == 0 {
//		rglog.SystemError("freeCache启动失败，未设置容量|" + memorySize)
//		return
//	}
//	freecacheUp(int(configSize))
//	go listen(redisClientType, redisClient)
//}
//
//// listen 订阅redis
//// @Param   :
//// @Return  :
//// @Author  : LiJunDong
//// @Time    : 2022-06-01
//func listen(clientType redisClientTypeEnum, clientInterface rgredis.RedisClientITF) {
//	keyPrefix := rgconfig.GetStr(redisKeyPrefix)
//	redisDB := rgconfig.GetStr(rgconst.ConfigKeyRedisDB)
//	if redisDB == "" {
//		redisDB = "0"
//	}
//	baseChanName := "__keyspace@" + redisDB + "__:" // <LiJunDong : 2022-06-01 17:41:18> --- 取配置文件中sys_redis_db, 默认为0
//	channelName := baseChanName
//	if keyPrefix != "" {
//		channelName += keyPrefix + "*"
//	}
//	var pubsub *redis.PubSub
//	if clientType == 1 {
//		redisLink, err := clientInterface.GetClient()
//		if err != nil {
//			rgo.This.Log.Error("单机redis链接转换失败")
//			return
//		}
//		client, ok := redisLink.(*redis.Client)
//		if !ok {
//			rgo.This.Log.Error("单机redis链接转换失败")
//			return
//		}
//		pubsub = client.PSubscribe(channelName)
//	}
//	if clientType == 2 {
//		redisLink, err := clientInterface.GetClient()
//		if err != nil {
//			rgo.This.Log.Error("redisCluster链接转换失败")
//			return
//		}
//		client, ok := redisLink.(*redis.ClusterClient)
//		if !ok {
//			rgo.This.Log.Error("redisCluster链接转换失败")
//			return
//		}
//		pubsub = client.PSubscribe(channelName)
//	}
//	ch := pubsub.Channel()
//	for msg := range ch {
//		// <LiJunDong : 2022-06-01 13:48:12> --- msg.String(): Message<__keyspace@0__:rgo:demo_cron:1: set>
//		// <LiJunDong : 2022-06-01 13:48:38> --- msg.Channel: __keyspace@0__:rgo:demo_cron:1
//		// <LiJunDong : 2022-06-01 13:48:38> --- msg.Payload: set
//		// <LiJunDong : 2022-06-01 13:48:38> --- msg.Pattern: msg.Payload
//		//log.Println("listen message", msg.String(), msg.Channel, msg.Payload, msg.Pattern)
//		key := msg.Channel[len(baseChanName):]
//		action := msg.Payload
//		del(action, key)
//	}
//}
//
//// del
//// @Param   : key string
//// @Return  :
//// @Author  : LiJunDong
//// @Time    : 2022-06-01
//func del(action string, keys ...string) {
//	for _, item := range keys {
//		Client.Del([]byte(item))
//	}
//}
//
//// freecacheUp 启动freecahche
//// @Param   :
//// @Return  :
//// @Author  : LiJunDong
//// @Time    : 2022-06-01
//func freecacheUp(size int) {
//	cacheSize := size * 1024 * 1024
//	Client = freecache.NewCache(cacheSize)
//}
//
//func GetRedisClient() (rgredis.RedisClientITF, error) {
//	return redisClient, nil
//}
