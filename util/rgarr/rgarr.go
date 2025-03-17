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
func Inter[T int | int64 | string](param1, param2 []T) (res []T) {
	defer func() {
		if err := recover(); err != nil {
			logStr, _ := json.Marshal([][]T{param1, param2})
			rglog.SystemError("recover捕获错误|" + fmt.Errorf("internal error: %v", err).Error() + "| data :" + string(logStr))
		}
	}()
	if len(param1) == 0 {
		return param1
	}
	if len(param2) == 0 {
		return param2
	}
	tmp := make(map[T]struct{}, len(param2)+len(param1))
	for _, item := range param2 {
		tmp[item] = struct{}{}
	}
	res = make([]T, 0, len(param1))
	for _, item := range param1 {
		if _, ok := tmp[item]; ok {
			res = append(res, item)
		}
	}
	return res
}

// Union 并集
// @Param   :
// @Return  :
// @Author  : LiJunDong
// @Time    : 2022-06-01
func Union[T int | int64 | string](param1, param2 []T) (res []T) {
	tmp := make(map[T]struct{}, len(param1)+len(param2))
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
func Diff[T int | int64 | string](param1, param2 []T) (res []T) {
	if len(param2) == 0 || len(param1) == 0 {
		return param1
	}
	tmp := make(map[T]struct{}, len(param2))
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
func Unique[T int | int64 | string](param []T) (data []T) {
	if len(param) == 0 {
		return param
	}
	tmp := make(map[T]struct{}, len(param))
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
func Reverse[T int | int64 | string](s []T) []T {
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

func StrInArray(val string, haystack []string) (exists bool) {
	for _, e := range haystack {
		if e == val {
			return true
		}
	}
	return false

}

func IntInArray(val int, haystack []int) (exists bool) {
	for _, e := range haystack {
		if e == val {
			return true
		}
	}
	return false
}

func InArray[T int | int64 | string](val T, haystack []T) (exists bool) {
	for _, e := range haystack {
		if e == val {
			return true
		}
	}
	return false
}

func IsEmptySlice[T int64 | string | float64](param []T) bool {
	if len(param) == 0 {
		return true
	}
	return false
}

func InSlice[T int | int64 | string](s []T, i T) bool {
	m := SliceToMap(s)
	_, ok := m[i]
	return ok
}

func SliceIntToStr[T int64 | int](param []T) (data []string) {
	data = make([]string, 0, len(param))
	for _, item := range param {
		data = append(data, strconv.Itoa(int(item)))
	}
	return data
}

func SliceStrToInt64(param []string) (data []int64) {
	data = make([]int64, 0, len(param))
	for _, item := range param {
		intItem, _ := strconv.Atoi(item)
		data = append(data, int64(intItem))
	}
	return data
}

func SliceToMap[T int | int64 | string](param []T) (data map[T]struct{}) {
	data = make(map[T]struct{}, len(param))
	for _, item := range param {
		data[item] = struct{}{}
	}
	return data
}

func ArrayPopItem[T int | int64 | string](item T, array []T) (newArray []T) {
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
func ArrayPopItems[T int | int64 | string](items []T, array []T) (newArray []T) {
	if len(array) == 0 {
		return array
	}
	newArray = array
	for _, item := range items {
		newArray = ArrayPopItem(item, newArray)
	}
	return newArray
}
