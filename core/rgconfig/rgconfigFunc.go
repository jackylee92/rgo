package rgconfig

/* 说明 <LiJunDong : 2022-03-02 14:05:42> --- 初始化获取配置
 *
 */

// <LiJunDong : 2022-03-02 16:08:36> --- 各种缓存标准interface
type ConfigInterface interface {
	Get(string) interface{}             // 获取interface类型
	GetStr(string) string               // 获取String类型
	GetInt(string) int64                // 获取int64类型
	GetBool(string) bool                // 获取bool类型
	GetStrMap(string) map[string]string // 获取map[string]string类型
	GetStrSlice(string) []string        // 获取[]string类型
	GetContent() string                 // 获取所有配置文件string
	Isset(string) bool                  // 判断配置是否存在
	Reload() error                      // 重新加载配置
	Load() func() error
}

/*
* @Content : 获取配置
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2022-03-01
 */
func Get(param string) (data interface{}) {
	if configPool == nil {
		return data
	}
	data = configPool.Get(param)
	return data
}

/*
* @Content : 获取配置
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2022-03-01
 */
func GetStr(param string) (data string) {
	if configPool == nil {
		return data
	}
	data = configPool.GetStr(param)
	return data
}

/*
* @Content : 获取配置
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2022-03-01
 */
func GetStrMap(param string) (data map[string]string) {
	if configPool == nil {
		return data
	}
	data = configPool.GetStrMap(param)
	return data
}

/*
* @Content : 获取配置
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2022-03-01
 */
func GetStrSlice(param string) (data []string) {
	if configPool == nil {
		return data
	}
	data = configPool.GetStrSlice(param)
	return data
}

/*
* @Content : 获取int64类型
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2022-05-11
 */
func GetInt(param string) (data int64) {
	if configPool == nil {
		return data
	}
	data = configPool.GetInt(param)
	return data
}

/*
* @Content : 获取配置
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2022-03-01
 */
func GetBool(param string) (data bool) {
	if configPool == nil {
		return data
	}
	data = configPool.GetBool(param)
	return data
}

/*
* @Content : 获取配置
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2022-03-01
 */
func GetContent() (data string) {
	if configPool == nil {
		return data
	}
	data = configPool.GetContent()
	return data
}

/*
* @Content : 是否存在
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2022-03-04
 */
func Isset(param string) (ok bool) {
	return ok
}

/*
* @Content : 重新加载
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2022-03-04
 */
func Reload() (err error) {
	return err
}
