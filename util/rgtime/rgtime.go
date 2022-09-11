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

// NowDateTime 获取当前时间 年月日时分秒
// @Param   :
// @Return  :
// @Author  : LiJunDong
// @Time    : 2022-05-28
func NowDateTime() string {
    return time.Now().Format(rgconst.GoTimeFormat)
}
