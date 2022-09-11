package rgredis

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/jackylee92/rgo/core/rgconfig"
	"github.com/jackylee92/rgo/core/rgglobal/rgconst"
	"github.com/jackylee92/rgo/core/rgglobal/rgerror"
	"github.com/jackylee92/rgo/core/rglog"

	"github.com/go-redis/redis"
)

// <LiJunDong : 2022-03-30 21:21:22> --- redis配置
type config struct {
	Typ          int      `json:"type"`           // 链接方式 1:单机链接 2:集群链接 3:读写分离
	Addr         string   `json:"addr"`           // 单机链接的地址
	Addrs        []string `json:"addrs"`          // 集群链接地址
	Password     string   `json:"password"`       // 链接认证
	PoolSize     int      `json:"pool_size"`      // 链接池大小
	MinIdleConns int      `json:"min_idle_conns"` // 初始化连接数
	DB           int      `json:"db"`             // DB
}

var globalClient *redis.Client
var globalClusterClient *redis.ClusterClient
var clientType int

// Start 启动
// @Param   :
// @Return  :
// @Author  : LiJunDong
// @Time    : 2022-05-31
func Start() {
	c, ok := getConfig()
	if !ok {
		return
	}
	addrLog := ""
	clientType = c.Typ
	switch c.Typ {
	case 1:
		// <LiJunDong : 2022-03-30 21:52:42> --- 单机
		globalClient = redis.NewClient(&redis.Options{
			Addr:         c.Addr,
			Password:     c.Password, // no password set
			DB:           c.DB,
			PoolSize:     c.PoolSize,     // 连接池最大socket连接数，默认为10倍CPU数， 10 * runtime.NumCPU
			MinIdleConns: c.MinIdleConns, //在启动阶段创建指定数量的Idle连接，并长期维持idle状态的连接数不少于指定数量；。
		})
		addrLog = c.Addr
		err := globalClient.Ping().Err()
		if err != nil {
			panic("redis链接失败")
		}
	case 2:
		// <LiJunDong : 2022-03-30 21:52:55> --- 集群
		globalClusterClient = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:        c.Addrs,
			Password:     c.Password,
			PoolSize:     c.PoolSize,     // 连接池最大socket连接数，默认为10倍CPU数， 10 * runtime.NumCPU
			MinIdleConns: c.MinIdleConns, //在启动阶段创建指定数量的Idle连接，并长期维持idle状态的连接数不少于指定数量；。
		})
		addrLog = strings.Join(c.Addrs, ",")
		err := globalClusterClient.Ping().Err()
		if err != nil {
			panic("redis链接失败｜" + err.Error())
		}
	case 3:
		// TODO <LiJunDong : 2022-03-30 21:53:04> --- 读写分离
		panic("暂不支持读写分离模式")
	}
	rglog.SystemInfo("启动项【redis】" + addrLog + ":成功")
}

/*
* @Content : 获取根据不同配置文件方式redis链接配置
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2022-05-12
 */
func getConfig() (data config, isset bool) {
	if rgconfig.Config.Config == "apollo" {
		redisJson := rgconfig.GetStr(rgconst.ConfigKeyRedis)
		if redisJson == "" {
			return data, false
		}
		err := json.Unmarshal([]byte(redisJson), &data)
		if err != nil {
			rglog.SystemError(rgerror.ErrorRedisConfigNil)
			panic(err)
		}
	}
	if rgconfig.Config.Config == "file" {
		typ := rgconfig.GetInt(rgconst.ConfigKeyRedisType)
		addr := rgconfig.GetStr(rgconst.ConfigKeyRedisAddr)
		addrs := rgconfig.GetStrSlice(rgconst.ConfigKeyRedisAddrs)
		password := rgconfig.GetStr(rgconst.ConfigKeyRedisPassword)
		poolSize := rgconfig.GetInt(rgconst.ConfigKeyRedisPoolSize)
		minIdleConns := rgconfig.GetInt(rgconst.ConfigKeyRedisMinIdleConns)
		db := rgconfig.GetInt(rgconst.ConfigKeyRedisDB)
		data = config{
			Typ:          int(typ),
			Addr:         addr,
			Addrs:        addrs,
			Password:     password,
			PoolSize:     int(poolSize),
			MinIdleConns: int(minIdleConns),
			DB:           int(db),
		}
		if data.Typ != 1 && data.Typ != 2 && data.Typ != 3 {
			return data, false
		}
		if data.Typ == 1 {
			if data.Addr == "" {
				return data, false
			}
		}
		if data.Typ == 2 {
			if len(data.Addrs) == 0 {
				return data, false
			}
		}
	}
	return data, true
}

/*
* @Content : 根据配置获取链接
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2022-03-30
 */
func GetClient() (redis.Cmdable, error) {
	switch clientType {
	case 1:
		return globalClient, nil
	case 2:
		return globalClusterClient, nil
	default:
		return nil, errors.New("未获取到redis链接对象")
	}
}

// GetClientType 获取使用的类型
// @Param   :
// @Return  : typ int
// @Author  : LiJunDong
// @Time    : 2022-06-01
func GetClientType() (typ int) {
	return clientType
}

/*
* @Content : string类型 获取
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2021-09-09
 */
func Get(key string) (string, error) {
	client, err := GetClient()
	if err != nil {
		return "", errors.New(rgerror.ErrorRedisClientNil)
	}
	data, err := client.Get(key).Result()
	return data, err
}

/*
* @Content : 设置可以有失效时间的string类型 t可以不传
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2021-09-09
 */
func Setex(key string, val string, t time.Duration) (bool, error) {
	if key == "" {
		err := errors.New("redis设置有失效时间的string类型结果失败")
		return false, err
	}
	if val == "" {
		err := errors.New("redis设置有失效时间的string类型结果失败,val为空")
		return false, err
	}
	client, err := GetClient()
	if err != nil {
		return false, errors.New(rgerror.ErrorRedisClientNil)
	}
	_, err = client.Set(key, val, t).Result()
	if err != nil {
		err := errors.New("redis设置有失效时间的string类型结果失败" + err.Error())
		return false, err
	}
	return true, nil
}

/*
* @Content : 设置不存在的有失效时间的string类型，如果存在则返回false
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2021-09-09
 */
func Setnx(key string, val string, t time.Duration) (bool, error) {
	if key == "" {
		err := errors.New("Redis设置不存在的key结果失败，key为空")
		return false, err
	}
	if val == "" {
		err := errors.New("Redis设置不存在的key结果失败，val为空")
		return false, err
	}
	client, err := GetClient()
	if err != nil {
		return false, errors.New(rgerror.ErrorRedisClientNil)
	}
	result, err := client.SetNX(key, val, t).Result()
	if err != nil {
		err := errors.New("Redis设置不存在的key结果失败" + err.Error())
		return false, err
	}
	return result, nil
}

/*
* @Content : 删除key
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2021-09-09
 */
func Del(key string) bool {
	if key == "" {
		return false
	}
	client, err := GetClient()
	if err != nil {
		return false
	}
	_, err = client.Del(key).Result()
	if err != nil {
		return false
	}
	return true
}

/*
* @Content : 自增
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2021-09-09
 */
func Incrby(key string, value int64) (int64, error) {
	client, err := GetClient()
	if err != nil {
		return 0, errors.New(rgerror.ErrorRedisClientNil)
	}
	return client.IncrBy(key, value).Result()
}

/*
* @Content : list头部弹出一个
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2021-09-09
 */
func LPop(key string) (string, error) {
	client, err := GetClient()
	if err != nil {
		return "", errors.New(rgerror.ErrorRedisClientNil)
	}
	return client.LPop(key).Result()
}

/*
* @Content : list尾部弹出一个
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2021-09-09
 */
func RPop(key string) (string, error) {
	client, err := GetClient()
	if err != nil {
		return "", errors.New(rgerror.ErrorRedisClientNil)
	}
	return client.RPop(key).Result()
}

/*
* @Content : list头部插入一个
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2021-09-09
 */
func LPush(key string, value string) (int64, error) {
	client, err := GetClient()
	if err != nil {
		return 0, errors.New(rgerror.ErrorRedisClientNil)
	}
	return client.LPush(key, value).Result()
}

/*
* @Content : list尾部插入一个
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2021-09-09
 */
func RPush(key string, value string) (int64, error) {
	client, err := GetClient()
	if err != nil {
		return 0, errors.New(rgerror.ErrorRedisClientNil)
	}
	return client.RPush(key, value).Result()
}

/*
* @Content : list获取所有
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2021-09-09
 */
func LRange(key string, start int64, end int64) ([]string, error) {
	client, err := GetClient()
	if err != nil {
		return []string{}, errors.New(rgerror.ErrorRedisClientNil)
	}
	return client.LRange(key, start, end).Result()
}

/*
* @Content : 过期设置
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2021-10-20
 */
func Expire(key string, t time.Duration) (bool, error) {
	client, err := GetClient()
	if err != nil {
		return false, errors.New(rgerror.ErrorRedisClientNil)
	}
	return client.Expire(key, t).Result()
}

/*
* @Content : hash
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2021-11-05
 */
func HSet(table, field string, value interface{}) (bool, error) {
	client, err := GetClient()
	if err != nil {
		return false, errors.New(rgerror.ErrorRedisClientNil)
	}
	return client.HSet(table, field, value).Result()
}

/*
* @Content : 批量获取
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2021-11-10
 */
func MGet(param ...string) ([]interface{}, error) {
	param = unique(param)
	client, err := GetClient()
	if err != nil {
		return []interface{}{}, errors.New(rgerror.ErrorRedisClientNil)
	}
	return client.MGet(param...).Result()
}

/*
* @Content : SAdd
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2021-11-11
 */
func SAdd(key string, param ...interface{}) (int64, error) {
	client, err := GetClient()
	if err != nil {
		return 0, errors.New(rgerror.ErrorRedisClientNil)
	}
	return client.SAdd(key, param...).Result()
}

/*
* @Content : SRem
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2021-11-11
 */
func SRem(key string, param ...interface{}) (int64, error) {
	client, err := GetClient()
	if err != nil {
		return 0, errors.New(rgerror.ErrorRedisClientNil)
	}
	return client.SRem(key, param...).Result()
}

/*
* @Content : SMembers
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2021-11-11
 */
func SMembers(key string) ([]string, error) {
	client, err := GetClient()
	if err != nil {
		return []string{}, errors.New(rgerror.ErrorRedisClientNil)
	}
	return client.SMembers(key).Result()
}

/*
* @Content : 集合取交集，暂且不使用redis做交集差集，数值过多可能会造成redis阻塞，相对于redis升级，服务器升级的扩展比较容易
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2021-11-11
 */
func SInter(param ...string) ([]string, error) {
	client, err := GetClient()
	if err != nil {
		return []string{}, errors.New(rgerror.ErrorRedisClientNil)
	}
	return client.SInter(param...).Result()
}

/*
* @Content : 集合取差集
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2021-11-11
 */
func SDiff(param ...string) ([]string, error) {
	client, err := GetClient()
	if err != nil {
		return []string{}, errors.New(rgerror.ErrorRedisClientNil)
	}
	return client.SDiff(param...).Result()
}

/*
* @Content : 并集
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2021-11-16
 */
func SUnion(param ...string) ([]string, error) {
	client, err := GetClient()
	if err != nil {
		return []string{}, errors.New(rgerror.ErrorRedisClientNil)
	}
	return client.SUnion(param...).Result()
}

/*
* @Content : eval
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2021-11-22
 */
func Eval(script string, keys []string, args []interface{}) (interface{}, error) {
	client, err := GetClient()
	if err != nil {
		return nil, errors.New(rgerror.ErrorRedisClientNil)
	}
	return client.Eval(script, keys, args...).Result()
}

/*
* @Content : hash字段incrby
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2021-11-26
 */
func HIncrby(key, field string, value int64) (int64, error) {
	client, err := GetClient()
	if err != nil {
		return 0, errors.New(rgerror.ErrorRedisClientNil)
	}
	return client.HIncrBy(key, field, value).Result()
}

/*
* @Content : HGet
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2021-11-29
 */
func HGet(key string, field string) (string, error) {
	client, err := GetClient()
	if err != nil {
		return "", errors.New(rgerror.ErrorRedisClientNil)
	}
	return client.HGet(key, field).Result()
}

/*
* @Content : HDel
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2022-02-18
 */
func HDel(key string, field string) (int64, error) {
	client, err := GetClient()
	if err != nil {
		return 0, errors.New(rgerror.ErrorRedisClientNil)
	}
	return client.HDel(key, field).Result()
}

/*
* @Content : HMGet
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2021-11-29
 */
func HMGet(key string, field []string) ([]interface{}, error) {
	field = unique(field)
	client, err := GetClient()
	if err != nil {
		return nil, errors.New(rgerror.ErrorRedisClientNil)
	}
	return client.HMGet(key, field...).Result()
}

/*
* @Content : HGetAll
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2021-11-29
 */
func HGetAll(key string) (map[string]string, error) {
	client, err := GetClient()
	if err != nil {
		return nil, errors.New(rgerror.ErrorRedisClientNil)
	}
	return client.HGetAll(key).Result()
}

/*
* @Content : HMget
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2021-11-29
 */
func HMSet(key string, fields map[string]interface{}) (string, error) {
	client, err := GetClient()
	if err != nil {
		return "", errors.New(rgerror.ErrorRedisClientNil)
	}
	return client.HMSet(key, fields).Result()
}

/*
* @Content : 数组去重
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2021-11-16
 */
func unique(param []string) (data []string) {
	if len(param) == 0 {
		return param
	}
	tmp := make(map[string]struct{}, len(param))
	cut := 0
	for i := 0; i < len(param); i++ {
		item := param[i]
		if _, ok := tmp[item]; ok {
			continue
		}
		param[cut] = item
		tmp[item] = struct{}{}
		cut = cut + 1
	}
	return param[:cut]
}

func ZrangeByScore(key string, min string, max string, offset int64, count int64) ([]string, error) {
	client, err := GetClient()
	if err != nil {
		return nil, errors.New(rgerror.ErrorRedisClientNil)
	}
	return client.ZRangeByScore(key, redis.ZRangeBy{Min: min, Max: max, Offset: offset, Count: count}).Result()
}

func Zrank(key string, member string) (int64, error) {
	client, err := GetClient()
	if err != nil {
		return 0, errors.New(rgerror.ErrorRedisClientNil)
	}
	return client.ZRank(key, member).Result()
}

func Zscore(key string, member string) (float64, error) {
	client, err := GetClient()
	if err != nil {
		return 0, errors.New(rgerror.ErrorRedisClientNil)
	}
	return client.ZScore(key, member).Result()
}

func Zadd(key string, score float64, member string) (int64, error) {
	client, err := GetClient()
	if err != nil {
		return 0, errors.New(rgerror.ErrorRedisClientNil)
	}
	return client.ZAdd(key, redis.Z{Member: member, Score: score}).Result()
}

func SetPersist(key string) (bool, error) {
	client, err := GetClient()
	if err != nil {
		return false, errors.New(rgerror.ErrorRedisClientNil)
	}
	return client.Persist(key).Result()
}

func HKeys(key string) ([]string, error) {
	var res []string
	client, err := GetClient()
	if err != nil {
		return res, errors.New(rgerror.ErrorRedisClientNil)
	}
	return client.HKeys(key).Result()
}
func SIsMember(key string, args interface{}) (bool, error) {
	client, err := GetClient()
	if err != nil {
		return false, errors.New(rgerror.ErrorRedisClientNil)
	}
	return client.SIsMember(key, args).Result()
}
