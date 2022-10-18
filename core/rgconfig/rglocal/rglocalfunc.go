package rglocal

import (
	"github.com/jackylee92/rgo/core/rgconfig"
	"github.com/jackylee92/rgo/core/rgjson"
)

// 实现config方法的结构体
type client struct{}

var config map[string]interface{}

/*
* @Content : 加载
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2022-03-08
 */
func (c client) Load() func() error {
	return c.load
}

/*
* @Content : 根据key查询
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2022-03-08
 */
func (c client) GetStr(param string) (data string) {
	dataTmp, ok := config[param]
	if !ok {
		return ""
	}
	return toString(dataTmp)
}

/*
* @Content : 根据key查询
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2022-03-08
 */
func (c client) GetStrMap(param string) (data map[string]string) {
	dataTmp, ok := config[param]
	if !ok {
		return data
	}
	dataStr, err := rgjson.Marshel(dataTmp)
	if err != nil {
		return data
	}
	rgjson.UnMarshel([]byte(dataStr), &data)
	return data
}

/*
* @Content : 根据key查询
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2022-03-08
 */
func (c client) GetStrSlice(param string) (data []string) {
	dataTmp, ok := config[param]
	if !ok {
		return data
	}
	dataStr, err := rgjson.Marshel(dataTmp)
	if err != nil {
		return data
	}
	rgjson.UnMarshel([]byte(dataStr), &data)
	return data
}

/*
* @Content : 获取int切片
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2022-03-10
 */
func (c client) GetInt(param string) (data int64) {
	dataTmp, ok := config[param]
	if !ok {
		return data
	}
	return toInt(dataTmp)
}

/*
* @Content : 获取bool类型配置
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2022-03-10
 */
func (c client) GetBool(param string) (data bool) {
	dataTmp, ok := config[param]
	if !ok {
		return data
	}
	return toBool(dataTmp)
}

/*
* @Content :
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2022-03-10
 */
func (c client) GetContent() (data string) {
	return string(rgconfig.ReadFile())
}

/*
* @Content : 根据key判断是否存在
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2022-03-08
 */
func (c client) Isset(param string) (data bool) {
	return data
}

/*
* @Content : 重新加载配置
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2022-03-08
 */
func (c client) Reload() (err error) {
	return err
}
