package rglocal

import (
	"strconv"
	"strings"

	"github.com/jackylee92/rgo/core/rgconfig"
	"gopkg.in/yaml.v2"
)

type LocalConfig struct{}
type CreateLocalConfig struct{}

const (
	Name string = "file"
)

/*
* @Content : init
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2022-03-02
 */
func Register() {
	rgconfig.Register(Name, client{})
}

/*
* @Content : 具体加载方法
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2022-05-11
 */
func (c client) load() (err error) {
	configByte := rgconfig.ReadFile()
	err = yaml.Unmarshal(configByte, &config)
	if err != nil {
		return err
	}
	return err
}

/*
* @Content : interface转string
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2022-05-11
 */
func toString(param interface{}) (data string) {
	switch param.(type) {
	case string:
		if value, ok := param.(string); ok {
			data = value
		}
	case float64:
		if value, ok := param.(float64); ok {
			data = strconv.FormatFloat(value, 'f', -1, 64)
		}
	case int:
		if value, ok := param.(int); ok {
			data = strconv.Itoa(value)
		}
	case int64:
		if value, ok := param.(int64); ok {
			data = strconv.FormatInt(value, 10)
		}
	case float32:
		if value, ok := param.(float32); ok {
			data = strconv.FormatFloat(float64(value), 'f', -1, 32)
		}
	case bool:
		if value, ok := param.(bool); ok {
			if value == true {
				data = "true"
			} else {
				data = "false"
			}
		}
	default:
		return ""
	}
	return data
}

/*
* @Content : interface转int64
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2022-05-11
 */
func toInt(param interface{}) (data int64) {
	switch param.(type) {
	case string:
		if value, ok := param.(string); ok {
			strToInt, _ := strconv.Atoi(value)
			data = int64(strToInt)
		}
	case float64:
		if value, ok := param.(float64); ok {
			data = int64(value)
		}
	case int:
		if value, ok := param.(int); ok {
			data = int64(value)
		}
	case int64:
		if value, ok := param.(int64); ok {
			data = value
		}
	case int32:
		if value, ok := param.(int32); ok {
			data = int64(value)
		}
	case float32:
		if value, ok := param.(float32); ok {
			data = int64(value)
		}
	case bool:
		if value, ok := param.(bool); ok {
			if value == true {
				data = 1
			} else {
				data = 0
			}
		}
	default:
		return 0
	}
	return data
}

/*
* @Content : interface转bool
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2022-05-11
 */
func toBool(param interface{}) (data bool) {
	dataTmp := toString(param)
	if "true" == strings.ToLower(dataTmp) {
		return true
	} else {
		return false
	}
}
