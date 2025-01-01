package rgredis

import (
	"context"
	"errors"
	"github.com/jackylee92/rgo/core/rgglobal/rgerror"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisClientITF interface {
	GetClient() (redis.Cmdable, error)
	Get(string) (string, error)
	Del(string) bool
	LRange(string, int64, int64) ([]string, error)
	MGet(...string) ([]interface{}, error)
	SMembers(string) ([]string, error)
	SInter(...string) ([]string, error)
	SDiff(...string) ([]string, error)
	SUnion(...string) ([]string, error)
	HGet(string, string) (string, error)
	HMGet(string, []string) ([]interface{}, error)
	HGetAll(string) (map[string]string, error)
}

func (client *Client) GetClient() (redis.Cmdable, error) {
	return client.linkObj, nil
}

// Get /*
func (client *Client) Get(ctx context.Context, key string) (string, error) {
	data, err := client.linkObj.Get(ctx, key).Result()
	return data, err
}

/*
* @Content : 设置可以有失效时间的string类型 t可以不传
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2021-09-09
 */
func (client *Client) Setex(ctx context.Context, key string, val string, t time.Duration) (bool, error) {
	if client.linkObj == nil {
		return false, errors.New(rgerror.ErrorRedisClientNil)
	}
	if key == "" {
		err := errors.New("redis设置有失效时间的string类型结果失败")
		return false, err
	}
	if val == "" {
		err := errors.New("redis设置有失效时间的string类型结果失败,val为空")
		return false, err
	}
	_, err := client.linkObj.Set(ctx, key, val, t).Result()
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
func (client *Client) Setnx(ctx context.Context, key string, val string, t time.Duration) (bool, error) {
	if client.linkObj == nil {
		return false, errors.New(rgerror.ErrorRedisClientNil)
	}
	if key == "" {
		err := errors.New("redis设置不存在的key结果失败，key为空")
		return false, err
	}
	if val == "" {
		err := errors.New("redis设置不存在的key结果失败，val为空")
		return false, err
	}

	result, err := client.linkObj.SetNX(ctx, key, val, t).Result()
	if err != nil {
		err := errors.New("redis设置不存在的key结果失败" + err.Error())
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
func (client *Client) Del(ctx context.Context, key string) bool {
	if client.linkObj == nil {
		return false
	}
	if key == "" {
		return false
	}
	_, err := client.linkObj.Del(ctx, key).Result()
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
func (client *Client) Incrby(ctx context.Context, key string, value int64) (int64, error) {
	if client.linkObj == nil {
		return 0, errors.New(rgerror.ErrorRedisClientNil)
	}
	return client.linkObj.IncrBy(ctx, key, value).Result()
}

/*
* @Content : list头部弹出一个
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2021-09-09
 */
func (client *Client) LPop(ctx context.Context, key string) (string, error) {
	if client.linkObj == nil {
		return "", errors.New(rgerror.ErrorRedisClientNil)
	}
	return client.linkObj.LPop(ctx, key).Result()
}

/*
* @Content : list尾部弹出一个
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2021-09-09
 */
func (client *Client) RPop(ctx context.Context, key string) (string, error) {
	if client.linkObj == nil {
		return "", errors.New(rgerror.ErrorRedisClientNil)
	}
	return client.linkObj.RPop(ctx, key).Result()
}

/*
* @Content : list头部插入一个
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2021-09-09
 */
func (client *Client) LPush(ctx context.Context, key string, value string) (int64, error) {
	if client.linkObj == nil {
		return 0, errors.New(rgerror.ErrorRedisClientNil)
	}
	return client.linkObj.LPush(ctx, key, value).Result()
}

/*
* @Content : list尾部插入一个
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2021-09-09
 */
func (client *Client) RPush(ctx context.Context, key string, value string) (int64, error) {
	if client.linkObj == nil {
		return 0, errors.New(rgerror.ErrorRedisClientNil)
	}
	return client.linkObj.RPush(ctx, key, value).Result()
}

/*
* @Content : list获取所有
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2021-09-09
 */
func (client *Client) LRange(ctx context.Context, key string, start int64, end int64) (res []string, err error) {
	if client.linkObj == nil {
		return res, errors.New(rgerror.ErrorRedisClientNil)
	}
	return client.linkObj.LRange(ctx, key, start, end).Result()
}

/*
* @Content : 过期设置
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2021-10-20
 */
func (client *Client) Expire(ctx context.Context, key string, t time.Duration) (bool, error) {
	if client.linkObj == nil {
		return false, errors.New(rgerror.ErrorRedisClientNil)
	}
	return client.linkObj.Expire(ctx, key, t).Result()
}

/*
* @Content : hash
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2021-11-05
 */
func (client *Client) HSet(ctx context.Context, table, field string, value interface{}) (int64, error) {
	if client.linkObj == nil {
		return 0, errors.New(rgerror.ErrorRedisClientNil)
	}
	return client.linkObj.HSet(ctx, table, field, value).Result()
}

/*
* @Content : 批量获取
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2021-11-10
 */
func (client *Client) MGet(ctx context.Context, param ...string) (res []interface{}, err error) {
	if client.linkObj == nil {
		return res, errors.New(rgerror.ErrorRedisClientNil)
	}
	param = unique(param)
	return client.linkObj.MGet(ctx, param...).Result()
}

/*
* @Content : SAdd
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2021-11-11
 */
func (client *Client) SAdd(ctx context.Context, key string, param ...interface{}) (res int64, err error) {
	if client.linkObj == nil {
		return res, errors.New(rgerror.ErrorRedisClientNil)
	}
	return client.linkObj.SAdd(ctx, key, param...).Result()
}

/*
* @Content : SRem
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2021-11-11
 */
func (client *Client) SRem(ctx context.Context, key string, param ...interface{}) (res int64, err error) {
	if client.linkObj == nil {
		return res, errors.New(rgerror.ErrorRedisClientNil)
	}
	return client.linkObj.SRem(ctx, key, param...).Result()
}

/*
* @Content : SMembers
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2021-11-11
 */
func (client *Client) SMembers(ctx context.Context, key string) (res []string, err error) {
	if client.linkObj == nil {
		return res, errors.New(rgerror.ErrorRedisClientNil)
	}
	return client.linkObj.SMembers(ctx, key).Result()
}

/*
* @Content : 集合取交集，暂且不使用redis做交集差集，数值过多可能会造成redis阻塞，相对于redis升级，服务器升级的扩展比较容易
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2021-11-11
 */
func (client *Client) SInter(ctx context.Context, param ...string) (res []string, err error) {
	if client.linkObj == nil {
		return res, errors.New(rgerror.ErrorRedisClientNil)
	}
	return client.linkObj.SInter(ctx, param...).Result()
}

/*
* @Content : 集合取差集
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2021-11-11
 */
func (client *Client) SDiff(ctx context.Context, param ...string) (res []string, err error) {
	if client.linkObj == nil {
		return res, errors.New(rgerror.ErrorRedisClientNil)
	}
	return client.linkObj.SDiff(ctx, param...).Result()
}

/*
* @Content : 并集
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2021-11-16
 */
func (client *Client) SUnion(ctx context.Context, param ...string) (res []string, err error) {
	if client.linkObj == nil {
		return res, errors.New(rgerror.ErrorRedisClientNil)
	}
	return client.linkObj.SUnion(ctx, param...).Result()
}

/*
* @Content : eval
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2021-11-22
 */
func (client *Client) Eval(ctx context.Context, script string, keys []string, args []interface{}) (res interface{}, err error) {
	if client.linkObj == nil {
		return res, errors.New(rgerror.ErrorRedisClientNil)
	}
	return client.linkObj.Eval(ctx, script, keys, args...).Result()
}

/*
* @Content : hash字段incrby
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2021-11-26
 */
func (client *Client) HIncrby(ctx context.Context, key, field string, value int64) (res int64, err error) {
	if client.linkObj == nil {
		return res, errors.New(rgerror.ErrorRedisClientNil)
	}
	return client.linkObj.HIncrBy(ctx, key, field, value).Result()
}

/*
* @Content : HGet
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2021-11-29
 */
func (client *Client) HGet(ctx context.Context, key string, field string) (res string, err error) {
	if client.linkObj == nil {
		return res, errors.New(rgerror.ErrorRedisClientNil)
	}
	return client.linkObj.HGet(ctx, key, field).Result()
}

/*
* @Content : HDel
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2022-02-18
 */
func (client *Client) HDel(ctx context.Context, key string, field string) (res int64, err error) {
	if client.linkObj == nil {
		return res, errors.New(rgerror.ErrorRedisClientNil)
	}
	return client.linkObj.HDel(ctx, key, field).Result()
}

/*
* @Content : HMGet
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2021-11-29
 */
func (client *Client) HMGet(ctx context.Context, key string, field []string) (res []interface{}, err error) {
	if client.linkObj == nil {
		return res, errors.New(rgerror.ErrorRedisClientNil)
	}
	field = unique(field)
	return client.linkObj.HMGet(ctx, key, field...).Result()
}

/*
* @Content : HGetAll
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2021-11-29
 */
func (client *Client) HGetAll(ctx context.Context, key string) (res map[string]string, err error) {
	if client.linkObj == nil {
		return res, errors.New(rgerror.ErrorRedisClientNil)
	}
	return client.linkObj.HGetAll(ctx, key).Result()
}

/*
* @Content : HMget
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2021-11-29
 */
func (client *Client) HMSet(ctx context.Context, key string, fields map[string]interface{}) (res bool, err error) {
	if client.linkObj == nil {
		return res, errors.New(rgerror.ErrorRedisClientNil)
	}
	return client.linkObj.HMSet(ctx, key, fields).Result()
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

func (client *Client) ZrangeByScore(ctx context.Context, key string, min string, max string, offset int64, count int64) (res []string, err error) {
	if client.linkObj == nil {
		return res, errors.New(rgerror.ErrorRedisClientNil)
	}
	return client.linkObj.ZRangeByScore(ctx, key, &redis.ZRangeBy{Min: min, Max: max, Offset: offset, Count: count}).Result()
}

func (client *Client) Zrank(ctx context.Context, key string, member string) (res int64, err error) {
	if client.linkObj == nil {
		return res, errors.New(rgerror.ErrorRedisClientNil)
	}
	return client.linkObj.ZRank(ctx, key, member).Result()
}

func (client *Client) Zscore(ctx context.Context, key string, member string) (res float64, err error) {
	if client.linkObj == nil {
		return res, errors.New(rgerror.ErrorRedisClientNil)
	}
	return client.linkObj.ZScore(ctx, key, member).Result()
}

func (client *Client) Zadd(ctx context.Context, key string, score float64, member string) (res int64, err error) {
	if client.linkObj == nil {
		return res, errors.New(rgerror.ErrorRedisClientNil)
	}
	return client.linkObj.ZAdd(ctx, key, redis.Z{Member: member, Score: score}).Result()
}

func (client *Client) SetPersist(ctx context.Context, key string) (res bool, err error) {
	if client.linkObj == nil {
		return res, errors.New(rgerror.ErrorRedisClientNil)
	}
	return client.linkObj.Persist(ctx, key).Result()
}

func (client *Client) HKeys(ctx context.Context, key string) (res []string, err error) {
	if client.linkObj == nil {
		return res, errors.New(rgerror.ErrorRedisClientNil)
	}
	return client.linkObj.HKeys(ctx, key).Result()
}

func (client *Client) SIsMember(ctx context.Context, key string, member interface{}) (res bool, err error) {
	if client.linkObj == nil {
		return res, errors.New(rgerror.ErrorRedisClientNil)
	}
	return client.linkObj.SIsMember(ctx, key, member).Result()
}
