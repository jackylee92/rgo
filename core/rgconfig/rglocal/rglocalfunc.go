package rglocal

import (
	"github.com/jackylee92/rgo/core/rgconfig"
	"github.com/jackylee92/rgo/core/rgjson"
)

// 实现config方法的结构体
type client struct{}

var config map[string]interface{}

func (c client) Load() func() error {
	return c.load
}

func (c client) Get(param string) (data interface{}) {
	dataTmp, ok := config[param]
	if !ok {
		return nil
	}
	return dataTmp
}

func (c client) GetStr(param string) (data string) {
	dataTmp, ok := config[param]
	if !ok {
		return ""
	}
	return toString(dataTmp)
}

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

func (c client) GetInt(param string) (data int64) {
	dataTmp, ok := config[param]
	if !ok {
		return data
	}
	return toInt(dataTmp)
}

func (c client) GetBool(param string) (data bool) {
	dataTmp, ok := config[param]
	if !ok {
		return data
	}
	return toBool(dataTmp)
}

func (c client) GetContent() (data string) {
	return string(rgconfig.ReadFile())
}

func (c client) Isset(param string) (data bool) {
	return data
}

func (c client) Reload() (err error) {
	return err
}
