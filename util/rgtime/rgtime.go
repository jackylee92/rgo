package rgtime

import (
	"github.com/jackylee92/rgo/core/rgglobal/rgconst"
	"time"
)

/*
 * @Content : rgtime
 * @Author  : LiJunDong
 * @Time    : 2022-05-28$
 */

func TimeToInt(param string) (t int64) {
	if param == "" {
		return 0
	}
	timeLayout := rgconst.GoTimeFormat     //转化所需模板
	loc, err := time.LoadLocation("Local") //重要：获取时区
	if err != nil {
		return t
	}
	theTime, err := time.ParseInLocation(timeLayout, param, loc) //使用模板在对应时区转化为time.time类型
	if err != nil {
		return t
	}
	t = theTime.Unix()
	return t
}

/*
* @Content : 当前时间戳
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2021-11-08
 */
func NowTimeInt() int64 {
	return time.Now().Unix()
}

// NowTimeMsInt 获取当前的毫秒级时间戳
func NowTimeMsInt() int64 {
	return time.Now().UnixNano() / 1e6
}

func NowDate() string {
	nowDate := time.Now().Format(rgconst.GoDateFormat)
	return nowDate
}

func NowTime() string {
	nowTime := time.Now().Format(rgconst.GoTimeFormat)
	return nowTime
}

func DateIntToDate(param int64, format string) string {
	tmp := time.Unix(param, 0)
	return tmp.Format(format)
}

func NowTimeChange(years, months, days int) (newTimeStr string) {
	t := time.Now()                           // 获取当前时间
	newTime := t.AddDate(years, months, days) // 前一个月的日期
	newTimeStr = DateIntToDate(newTime.Unix(), rgconst.GoTimeFormat)
	return newTimeStr
}
