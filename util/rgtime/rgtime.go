package rgtime

import "time"

/*
 * @Content : rgtime
 * @Author  : LiJunDong
 * @Time    : 2022-05-28$
 */

// NowDateTime 获取当前时间 年月日时分秒
// @Param   :
// @Return  :
// @Author  : LiJunDong
// @Time    : 2022-05-28
func NowDateTime() string {
    return time.Now().Format("2006-01-02 15:04:05")
}
