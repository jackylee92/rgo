package rgapollo

// 实现config方法的结构体
type client struct{}

func (c client) Load() func() error {
	return c.load
}

func (c client) Get(param string) (data interface{}) {
	data = configPool.GetValue(param)
	return data
}

func (c client) GetStr(param string) (data string) {
	data = configPool.GetValue(param)
	return data
}

func (c client) GetStrMap(param string) (data map[string]string) {
	return data
}

func (c client) GetStrSlice(param string) (data []string) {
	return data
}

func (c client) GetInt(param string) (data int64) {
	dataTmp := configPool.GetIntValue(param, 0)
	return int64(dataTmp)
}

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

func (c client) GetContent() (data string) {
	data = configPool.GetContent()
	return data
}

func (c client) Isset(param string) (res bool) {
	return res
}

func (c client) Reload() (err error) {
	return err
}
