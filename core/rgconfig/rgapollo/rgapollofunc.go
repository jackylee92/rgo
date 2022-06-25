package rgapollo

// 实现config方法的结构体
type client struct{}

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
	data = configPool.GetValue(param)
	return data
}

/*
* @Content : 根据key查询
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2022-03-08
 */
func (c client) GetStrMap(param string) (data map[string]string) {
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
	dataTmp := configPool.GetIntValue(param, 0)
	return int64(dataTmp)
}

/*
* @Content : 获取bool类型配置
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2022-03-10
 */
func (c client) GetBool(param string) (data bool) {
	// data = configPool.GetBoolValue(param, false)
	// <LiJunDong : 2022-04-11 13:19:29> --- 发现这个getBool获取到的是字符串类型，无法转为boo值，导致获取失败返回默认值false
	dataStr := configPool.GetValue(param)
	if dataStr == "true" {
		data = true
	} else {
		data = false
	}
	return data
}

/*
* @Content :
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2022-03-10
 */
func (c client) GetContent() (data string) {
	data = configPool.GetContent()
	return data
}

/*
* @Content : 根据key判断是否存在
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2022-03-08
 */
func (c client) Isset(param string) (res bool) {
	return res
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
