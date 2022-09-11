package rgcache

import (
	"errors"
	"github.com/jackylee92/rgo/core/rgjson"
	"github.com/jackylee92/rgo/core/rgmodel/rgredis"
	"github.com/jackylee92/rgo/util/rgarr"
)

var errorNoStart = errors.New("本地缓存freecache未启动")

/*
 * @Content : rgcache-实现多级缓存的组建，流程：本地内存->redis
 * 对外支持只查找
 * 添加修改数据绑定在查询之后
 * <LiJunDong : 2022-05-30 23:56:21> --- 查询数据的时候如何确保本地的数据和redis相同，也就是redis数据修改了本地这条数据未修改，此时查询到的数据就是错误的？
 * <LiJunDong : 2022-05-31 10:39:44> --- 通过redis订阅机制实现，注意以下：
 * 1. 所有项目使用的redis key都需要有一个统一的前缀，订阅会订阅这个前端的key的事件，时间包过：del、set、expire、expired、lpush等涉及到数据变化的指令，不区别数据类型
 * 2. 当list没有数据是会订阅到一条del时间
 * 3. redis需要执行config set notify-keyspace-events KEA 或者redis.conf中配置notify-keyspace-events KEA
 * 依赖组建 freecache、go-redis
 * @Author  : LiJunDong
 * @Time    : 2022-05-30$
 */

// Get
// @Param   :
// @Return  :
// @Author  : LiJunDong
// @Time    : 2022-05-31
func Get(key string) (data string, err error) {
	if Client == nil {
		return data, errorNoStart
	}
	cacheDataBt, err := Client.Get([]byte(key))
	if err != nil || len(cacheDataBt) == 0 { // <LiJunDong : 2022-06-01 14:48:33> --- 不存在
		data, _ = rgredis.Get(key)
		saveCache(key, data)
	} else {
		data = string(cacheDataBt)
	}
	return data, err
}

// LRange 查询集合
// @Param   :
// @Return  :
// @Author  : LiJunDong
// @Time    : 2022-05-31
func LRange(key string, start int64, end int64) (data []string, err error) {
	if Client == nil {
		return data, errorNoStart
	}
	cacheDataBt, err := Client.Get([]byte(key))
	if err != nil || len(cacheDataBt) == 0 { // <LiJunDong : 2022-06-01 14:48:33> --- 不存在
		allData, _ := rgredis.LRange(key, 0, -1) // <LiJunDong : 2022-06-01 15:00:03> --- 获取所有
		if len(allData) != 0 {
			cacheDataStr, _ := rgjson.Marshel(allData)
			cacheDataBt = []byte(cacheDataStr)
			saveCache(key, cacheDataStr)
		} else {
			return data, err
		}
	}
	var allData []string
	rgjson.UnMarshel(cacheDataBt, &allData)
	if int(start) > len(allData) || int(end) > len(allData) {
		return data, errors.New("数组越界")
	}
	data = allData[start:end]
	return data, err
}

// MGet 批量获取
// @Param   :
// @Return  :
// @Author  : LiJunDong
// @Time    : 2022-05-31
func MGet(param ...string) (data []interface{}, err error) {
	if Client == nil {
		return data, errorNoStart
	}
	for _, key := range param {
		cacheDataBt, err := Client.Get([]byte(key))
		if err != nil || len(cacheDataBt) == 0 { // <LiJunDong : 2022-06-01 14:48:33> --- 不存在
			dataItem, _ := rgredis.Get(key)
			if dataItem != "" {
				saveCache(key, dataItem)
			}
			data = append(data, dataItem)
		} else {
			data = append(data, string(cacheDataBt))
		}
	}
	return data, err
}

// SMembers 查询集合成员
// @Param   :
// @Return  :
// @Author  : LiJunDong
// @Time    : 2022-05-31
func SMembers(key string) (data []string, err error) {
	if Client == nil {
		return data, errorNoStart
	}
	cacheDataBt, err := Client.Get([]byte(key))
	if err != nil || len(cacheDataBt) == 0 { // <LiJunDong : 2022-06-01 14:48:33> --- 不存在
		data, _ = rgredis.SMembers(key)
		if len(data) != 0 {
			cacheDataStr, _ := rgjson.Marshel(data)
			saveCache(key, cacheDataStr)
		}
	} else {
		rgjson.UnMarshel(cacheDataBt, &data)
	}
	return data, err
}

// SInter 取交集
// @Param   :
// @Return  :
// @Author  : LiJunDong
// @Time    : 2022-05-31
func SInter(param ...string) (data []string, err error) {
	if Client == nil {
		return data, errorNoStart
	}
	for key, item := range param {
		var dataItem []string
		cacheDataBt, err := Client.Get([]byte(item))
		if err != nil || len(cacheDataBt) == 0 { // <LiJunDong : 2022-06-01 14:48:33> --- 不存在
			dataItem, _ = rgredis.SMembers(item)
			if len(dataItem) != 0 {
				cacheDataStr, _ := rgjson.Marshel(dataItem)
				saveCache(item, cacheDataStr)
			}
		} else {
			rgjson.UnMarshel(cacheDataBt, &dataItem)
		}
		if key == 0 {
			data = dataItem
		} else {
			data = rgarr.Inter(data, dataItem)
		}
		if len(data) == 0 {
			return data, err
		}
	}
	data = rgarr.Unique(data)
	return data, err
}

// SDiff 取差集
// @Param   :
// @Return  :
// @Author  : LiJunDong
// @Time    : 2022-05-31
func SDiff(param ...string) (data []string, err error) {
	if Client == nil {
		return data, errorNoStart
	}
	for _, item := range param {
		var dataItem []string
		cacheDataBt, err := Client.Get([]byte(item))
		if err != nil || len(cacheDataBt) == 0 { // <LiJunDong : 2022-06-01 14:48:33> --- 不存在
			dataItem, _ = rgredis.SMembers(item)
			if len(dataItem) != 0 {
				cacheDataStr, _ := rgjson.Marshel(dataItem)
				saveCache(item, cacheDataStr)
			}
		} else {
			rgjson.UnMarshel(cacheDataBt, &dataItem)
		}
		data = rgarr.Diff(data, dataItem)
	}
	data = rgarr.Unique(data)
	return data, err
}

// SUnion 取并集
// @Param   :
// @Return  :
// @Author  : LiJunDong
// @Time    : 2022-05-31
func SUnion(param ...string) (data []string, err error) {
	if Client == nil {
		return data, errorNoStart
	}
	for _, item := range param {
		var dataItem []string
		cacheDataBt, err := Client.Get([]byte(item))
		if err != nil || len(cacheDataBt) == 0 { // <LiJunDong : 2022-06-01 14:48:33> --- 不存在
			dataItem, _ = rgredis.SMembers(item)
			if len(dataItem) != 0 {
				cacheDataStr, _ := rgjson.Marshel(dataItem)
				saveCache(item, cacheDataStr)
			}
		} else {
			rgjson.UnMarshel(cacheDataBt, &dataItem)
		}
		data = append(data, dataItem...)
	}
	data = rgarr.Unique(data)
	return data, err
}

// HGet hash 查询
// @Param   :
// @Return  :
// @Author  : LiJunDong
// @Time    : 2022-05-31
func HGet(key string, field string) (data string, err error) {
	if Client == nil {
		return data, errorNoStart
	}
	hashData := make(map[string]string)
	cacheDataBt, err := Client.Get([]byte(key))
	if err != nil || len(cacheDataBt) == 0 { // <LiJunDong : 2022-06-01 14:48:33> --- 不存在
		hashData, _ = rgredis.HGetAll(key)
		if len(hashData) != 0 {
			cacheDataStr, _ := rgjson.Marshel(hashData)
			saveCache(key, cacheDataStr)
		}
	} else {
		rgjson.UnMarshel(cacheDataBt, &hashData)
	}
	data = hashData[field]
	return data, err
}

// HMGet hash 批量获取
// @Param   :
// @Return  :
// @Author  : LiJunDong
// @Time    : 2022-05-31
func HMGet(key string, field []string) (data []interface{}, err error) {
	if Client == nil {
		return data, errorNoStart
	}
	hashData := make(map[string]string)
	cacheDataBt, err := Client.Get([]byte(key))
	if err != nil || len(cacheDataBt) == 0 { // <LiJunDong : 2022-06-01 14:48:33> --- 不存在
		hashData, _ = rgredis.HGetAll(key)
		if len(hashData) != 0 {
			cacheDataStr, _ := rgjson.Marshel(hashData)
			saveCache(key, cacheDataStr)
		}
	} else {
		rgjson.UnMarshel(cacheDataBt, &hashData)
	}
	for _, item := range field {
		data = append(data, hashData[item])
	}
	return data, err
}

// HGetAll hash获取所有
// @Param   :
// @Return  :
// @Author  : LiJunDong
// @Time    : 2022-05-31
func HGetAll(key string) (data map[string]string, err error) {
	if Client == nil {
		return data, errorNoStart
	}
	cacheDataBt, err := Client.Get([]byte(key))
	if err != nil || len(cacheDataBt) == 0 { // <LiJunDong : 2022-06-01 14:48:33> --- 不存在
		data, _ = rgredis.HGetAll(key)
		if len(data) != 0 {
			cacheDataStr, _ := rgjson.Marshel(data)
			saveCache(key, cacheDataStr)
		}
	} else {
		rgjson.UnMarshel(cacheDataBt, &data)
	}
	return data, err
}

// saveCache
// @Param   : key string, data string
// @Return  :
// @Author  : LiJunDong
// @Time    : 2022-06-01
func saveCache(key string, data string) {
	if data == "" || key == "" {
		return
	}
	Client.Set([]byte(key), []byte(data), 0)
}
