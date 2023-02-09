package rgarr

import (
	"encoding/json"
	"fmt"
	"github.com/jackylee92/rgo/core/rglog"
	"sort"
	"strconv"
)

// Inter 交集 即存在param1又存在param2，不去重
// @Param   :
// @Return  :
// @Author  : LiJunDong
// @Time    : 2022-06-01
func Inter(param1, param2 []string) (res []string) {
	defer func() {
		if err := recover(); err != nil {
			logStr, _ := json.Marshal([][]string{param1, param2})
			rglog.SystemError("recover捕获错误|" + fmt.Errorf("internal error: %v", err).Error() + "| data :" + string(logStr))
		}
	}()
	if len(param1) == 0 {
		return param1
	}
	if len(param2) == 0 {
		return param2
	}
	tmp := make(map[string]struct{}, len(param2)+len(param1))
	for _, item := range param2 {
		tmp[item] = struct{}{}
	}
	cut := 0
	for _, item := range param1 {
		if _, ok := tmp[item]; ok {
			param1[cut] = item
			cut = cut + 1
		}
	}
	return param1[:cut]
}

// Union 并集
// @Param   :
// @Return  :
// @Author  : LiJunDong
// @Time    : 2022-06-01
func Union(param1, param2 []string) (res []string) {
	tmp := make(map[string]struct{}, len(param1)+len(param2))
	for _, item := range param1 {
		tmp[item] = struct{}{}
	}
	for _, item := range param2 {
		tmp[item] = struct{}{}
	}
	for key := range tmp {
		res = append(res, key)
	}
	return res
}

// Diff 差集 param1不存在param2中的元素集合
// @Param   :
// @Return  :
// @Author  : LiJunDong
// @Time    : 2022-06-01
func Diff(param1, param2 []string) (res []string) {
	if len(param2) == 0 || len(param1) == 0 {
		return param1
	}
	tmp := make(map[string]struct{}, len(param2))
	for _, item := range param2 {
		tmp[item] = struct{}{}
	}
	cut := 0
	for _, item := range param1 {
		inner := false
		if _, ok := tmp[item]; ok {
			inner = true
		}
		if inner {
			continue
		}
		param1[cut] = item
		cut = cut + 1
	}
	return param1[:cut]
}

// Unique 数组去重
// @Param   :
// @Return  :
// @Author  : LiJunDong
// @Time    : 2022-06-01
func Unique(param []string) (data []string) {
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

// Reverse 数组倒叙
// @Param   :
// @Return  :
// @Author  : LiJunDong
// @Time    : 2022-06-01
func Reverse(s []string) []string {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

// SortStringNumber 字符串数据排序
// @Param   :
// @Return  :
// @Author  : LiJunDong
// @Time    : 2022-06-01
func SortStringNumber(ss []string) []string {
	s := make([]int, 0, len(ss))
	for _, item := range ss {
		itemInt, _ := strconv.Atoi(item)
		s = append(s, itemInt)
	}
	sort.Ints(s)
	for key, item := range s {
		ss[key] = strconv.Itoa(item)
	}
	return ss
}

func InArrayStr(val string, haystack []string) (exists bool) {
	for _, e := range haystack {
		if e == val {
			return true
		}
	}
	return false

}

func InArrayInt(val int, haystack []int) (exists bool) {
	for _, e := range haystack {
		if e == val {
			return true
		}
	}
	return false
}

func InArrayInt64(val int64, haystack []int64) (exists bool) {
	for _, e := range haystack {
		if e == val {
			return true
		}
	}
	return false
}

func ArrayPopItem(item int64, array []int64) (newArray []int64) {
	if len(array) == 0 {
		return array
	}
	for key, val := range array {
		if val == item {
			return append(array[:key], array[(key+1):]...)
		}
	}
	return array
}
func ArrayPopItems(items []int64, array []int64) (newArray []int64) {
	if len(array) == 0 {
		return array
	}
	newArray = array
	for _, item := range items {
		newArray = ArrayPopItem(item, newArray)
	}
	return newArray
}

func KeyStringMapExist(key string, mapData map[string]string) (res bool) {
	if _, ok := mapData[key]; ok {
		return true
	} else {
		return false
	}
}

func KeyInt64MapExist(key int64, mapData map[int64]int64) (res bool) {
	if _, ok := mapData[key]; ok {
		return true
	} else {
		return false
	}
}

func Int64SliceToStringSlice(param []int64) (data []string) {
	if len(param) == 0 {
		return data
	}
	data = make([]string, 0, len(param))
	for _, item := range param {
		data = append(data, strconv.Itoa(int(item)))
	}
	return data
}

func StringSliceToInt64Slice(param []string) (data []int64) {
	if len(param) == 0 {
		return data
	}
	data = make([]int64, 0, len(param))
	for _, item := range param {
		value, err := strconv.Atoi(item)
		if err != nil {
			continue
		}
		data = append(data, int64(value))
	}
	return data
}
