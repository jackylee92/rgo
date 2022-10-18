package rgfloat

import "math"

/*
 * @Content : rgfloat
 * @Author  : LiJunDong
 * @Time    : 2022-10-19$
 */

// 浮点型 四舍五入转int64
func Round(x float64) int64 {
	return int64(math.Floor(x + 0.5))
}
