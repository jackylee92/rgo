package rgrediscluster

import (
	"errors"
	"github.com/jackylee92/rgo/core/rgglobal/rgerror"
	"time"

	"github.com/go-redis/redis"
)

/*Get
* @Content : string类型 获取
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2021-09-09
 */
func (client *Client) Get(key string) (string, error) {
	data, err := client.linkObj.Get(key).Result()
	return data, err
}

/*Setex
* @Content : 设置可以有失效时间的string类型 t可以不传
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2021-09-09
 */
func (client *Client) Setex(key string, val string, t time.Duration) (bool, error) {
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
	_, err := client.linkObj.Set(key, val, t).Result()
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
func (client *Client) Setnx(key string, val string, t time.Duration) (bool, error) {
	if client.linkObj == nil {
		return false, errors.New(rgerror.ErrorRedisClientNil)
	}
	if key == "" {
		err := errors.New("Redis设置不存在的key结果失败，key为空")
		return false, err
	}
	if val == "" {
		err := errors.New("Redis设置不存在的key结果失败，val为空")
		return false, err
	}

	result, err := client.linkObj.SetNX(key, val, t).Result()
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
func (client *Client) Del(key string) bool {
	if client.linkObj == nil {
		return false
	}
	if key == "" {
		return false
	}
	_, err := client.linkObj.Del(key).Result()
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
func (client *Client) Incrby(key string, value int64) (int64, error) {
	if client.linkObj == nil {
		return 0, errors.New(rgerror.ErrorRedisClientNil)
	}
	return client.linkObj.IncrBy(key, value).Result()
}

/*
* @Content : list头部弹出一个
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2021-09-09
 */
func (client *Client) LPop(key string) (string, error) {
	if client.linkObj == nil {
		return "", errors.New(rgerror.ErrorRedisClientNil)
	}
	return client.linkObj.LPop(key).Result()
}

/*
* @Content : list尾部弹出一个
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2021-09-09
 */
func (client *Client) RPop(key string) (string, error) {
	if client.linkObj == nil {
		return "", errors.New(rgerror.ErrorRedisClientNil)
	}
	return client.linkObj.RPop(key).Result()
}

/*
* @Content : list头部插入一个
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2021-09-09
 */
func (client *Client) LPush(key string, value string) (int64, error) {
	if client.linkObj == nil {
		return 0, errors.New(rgerror.ErrorRedisClientNil)
	}
	return client.linkObj.LPush(key, value).Result()
}

/*
* @Content : list尾部插入一个
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2021-09-09
 */
func (client *Client) RPush(key string, value string) (int64, error) {
	if client.linkObj == nil {
		return 0, errors.New(rgerror.ErrorRedisClientNil)
	}
	return client.linkObj.RPush(key, value).Result()
}

/*
* @Content : list获取所有
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2021-09-09
 */
func (client *Client) LRange(key string, start int64, end int64) (res []string, err error) {
	if client.linkObj == nil {
		return res, errors.New(rgerror.ErrorRedisClientNil)
	}
	return client.linkObj.LRange(key, start, end).Result()
}

/*
* @Content : 过期设置
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2021-10-20
 */
func (client *Client) Expire(key string, t time.Duration) (bool, error) {
	if client.linkObj == nil {
		return false, errors.New(rgerror.ErrorRedisClientNil)
	}
	return client.linkObj.Expire(key, t).Result()
}

/*
* @Content : hash
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2021-11-05
 */
func (client *Client) HSet(table, field string, value interface{}) (bool, error) {
	if client.linkObj == nil {
		return false, errors.New(rgerror.ErrorRedisClientNil)
	}
	return client.linkObj.HSet(table, field, value).Result()
}

/*
* @Content : SAdd
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2021-11-11
 */
func (client *Client) SAdd(key string, param ...interface{}) (res int64, err error) {
	if client.linkObj == nil {
		return res, errors.New(rgerror.ErrorRedisClientNil)
	}
	return client.linkObj.SAdd(key, param...).Result()
}

/*
* @Content : SRem
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2021-11-11
 */
func (client *Client) SRem(key string, param ...interface{}) (res int64, err error) {
	if client.linkObj == nil {
		return res, errors.New(rgerror.ErrorRedisClientNil)
	}
	return client.linkObj.SRem(key, param...).Result()
}

/*
* @Content : SMembers
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2021-11-11
 */
func (client *Client) SMembers(key string) (res []string, err error) {
	if client.linkObj == nil {
		return res, errors.New(rgerror.ErrorRedisClientNil)
	}
	return client.linkObj.SMembers(key).Result()
}

/*
* @Content : hash字段incrby
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2021-11-26
 */
func (client *Client) HIncrby(key, field string, value int64) (res int64, err error) {
	if client.linkObj == nil {
		return res, errors.New(rgerror.ErrorRedisClientNil)
	}
	return client.linkObj.HIncrBy(key, field, value).Result()
}

/*
* @Content : HGet
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2021-11-29
 */
func (client *Client) HGet(key string, field string) (res string, err error) {
	if client.linkObj == nil {
		return res, errors.New(rgerror.ErrorRedisClientNil)
	}
	return client.linkObj.HGet(key, field).Result()
}

/*
* @Content : HDel
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2022-02-18
 */
func (client *Client) HDel(key string, field string) (res int64, err error) {
	if client.linkObj == nil {
		return res, errors.New(rgerror.ErrorRedisClientNil)
	}
	return client.linkObj.HDel(key, field).Result()
}

/*
* @Content : HMGet
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2021-11-29
 */
func (client *Client) HMGet(key string, field []string) (res []interface{}, err error) {
	if client.linkObj == nil {
		return res, errors.New(rgerror.ErrorRedisClientNil)
	}
	field = unique(field)
	return client.linkObj.HMGet(key, field...).Result()
}

/*
* @Content : HGetAll
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2021-11-29
 */
func (client *Client) HGetAll(key string) (res map[string]string, err error) {
	if client.linkObj == nil {
		return res, errors.New(rgerror.ErrorRedisClientNil)
	}
	return client.linkObj.HGetAll(key).Result()
}

/*
* @Content : HMget
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2021-11-29
 */
func (client *Client) HMSet(key string, fields map[string]interface{}) (res string, err error) {
	if client.linkObj == nil {
		return res, errors.New(rgerror.ErrorRedisClientNil)
	}
	return client.linkObj.HMSet(key, fields).Result()
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

func (client *Client) ZrangeByScore(key string, min string, max string, offset int64, count int64) (res []string, err error) {
	if client.linkObj == nil {
		return res, errors.New(rgerror.ErrorRedisClientNil)
	}
	return client.linkObj.ZRangeByScore(key, redis.ZRangeBy{Min: min, Max: max, Offset: offset, Count: count}).Result()
}

func (client *Client) Zrank(key string, member string) (res int64, err error) {
	if client.linkObj == nil {
		return res, errors.New(rgerror.ErrorRedisClientNil)
	}
	return client.linkObj.ZRank(key, member).Result()
}

func (client *Client) Zscore(key string, member string) (res float64, err error) {
	if client.linkObj == nil {
		return res, errors.New(rgerror.ErrorRedisClientNil)
	}
	return client.linkObj.ZScore(key, member).Result()
}

func (client *Client) Zadd(key string, score float64, member string) (res int64, err error) {
	if client.linkObj == nil {
		return res, errors.New(rgerror.ErrorRedisClientNil)
	}
	return client.linkObj.ZAdd(key, redis.Z{Member: member, Score: score}).Result()
}

func (client *Client) SetPersist(key string) (res bool, err error) {
	if client.linkObj == nil {
		return res, errors.New(rgerror.ErrorRedisClientNil)
	}
	return client.linkObj.Persist(key).Result()
}

func (client *Client) HKeys(key string) (res []string, err error) {
	if client.linkObj == nil {
		return res, errors.New(rgerror.ErrorRedisClientNil)
	}
	return client.linkObj.HKeys(key).Result()
}

func (client *Client) SIsMember(key string, member interface{}) (res bool, err error) {
	if client.linkObj == nil {
		return res, errors.New(rgerror.ErrorRedisClientNil)
	}
	return client.linkObj.SIsMember(key, member).Result()
}
