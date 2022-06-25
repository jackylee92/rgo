package rgmath

import (
	"errors"
	"github.com/shopspring/decimal"
)

/*
 * @Content : rgmath
 * @Author  : LiJunDong
 * @Time    : 2022-06-18$
 */

// Add
// @Param   : param1, param2 interface{} 支持int， float, string
// @Return  : param decimal.Decimal
// @Author  : LiJunDong
// @Time    : 2022-06-18
func Add(param1, param2 interface{}) (data decimal.Decimal, err error) {
	p1, err := interfaceToDecimal(param1)
	if err != nil {
		return data, errors.New("Add失败，param1：" + err.Error())
	}
	p2, err := interfaceToDecimal(param2)
	if err != nil {
		return data, errors.New("Add失败，param2：" + err.Error())
	}
	data = p1.Add(p2)
    return data, err
}

// Sub
// @Param   : param1, param2 interface{}
// @Return  : data decimal.Decimal, err error
// @Author  : LiJunDong
// @Time    : 2022-06-18
func Sub(param1, param2 interface{}) (data decimal.Decimal, err error) {
	p1, err := interfaceToDecimal(param1)
	if err != nil {
		return data, errors.New("Sub失败，param1：" + err.Error())
	}
	p2, err := interfaceToDecimal(param2)
	if err != nil {
		return data, errors.New("Sub失败，param2：" + err.Error())
	}
	data = p1.Sub(p2)
	return data, err
}

// Mul
// @Param   : param1, param2 interface{}
// @Return  : data decimal.Decimal, err error
// @Author  : LiJunDong
// @Time    : 2022-06-18
func Mul(param1, param2 interface{}) (data decimal.Decimal, err error) {
	p1, err := interfaceToDecimal(param1)
	if err != nil {
		return data, errors.New("Mul失败，param1：" + err.Error())
	}
	p2, err := interfaceToDecimal(param2)
	if err != nil {
		return data, errors.New("Mul失败，param2：" + err.Error())
	}
	data = p1.Mul(p2)
	return data, err
}

// Div
// @Param   : param1, param2 interface{}
// @Return  : data decimal.Decimal, err error
// @Author  : LiJunDong
// @Time    : 2022-06-18
func Div(param1, param2 interface{}) (data decimal.Decimal, err error) {
	p1, err := interfaceToDecimal(param1)
	if err != nil {
		return data, errors.New("Div失败，param1：" + err.Error())
	}
	p2, err := interfaceToDecimal(param2)
	if err != nil {
		return data, errors.New("Div失败，Param2：" + err.Error())
	}
	if p2.String() == "0" {
		return data, errors.New("Div失败，Param2：被除数不能为0")
	}
	data = p1.Div(p2)
	return data, err
}

// interfaceToDecimal
// @Param   : param interface{}
// @Return  : data decimal.Decimal, err error
// @Author  : LiJunDong
// @Time    : 2022-06-18
func interfaceToDecimal(param interface{}) (data decimal.Decimal, err error) {
	switch param.(type) {
	case string:
		if value, ok := param.(string); ok {
			data, err = decimal.NewFromString(value)
		}else{
			err = errors.New("string类型转换decimal失败")
		}
	case float64:
		if value, ok := param.(float64); ok {
			data = decimal.NewFromFloat(value)
		}else{
			err = errors.New("float64类型转换decimal失败")
		}
	case int:
		if value, ok := param.(int); ok {
		data = decimal.NewFromInt(int64(value))
	}else{
		err = errors.New("int类型转换decimal失败")
	}
case int64:
		if value, ok := param.(int64); ok {
	data = decimal.NewFromInt(value)
	}else{
		err = errors.New("int64类型转换decimal失败")
	}
case int32:
		if value, ok := param.(int32); ok {
	data = decimal.NewFromInt32(value)
	}else{
		err = errors.New("int32类型转换decimal失败")
	}
	case float32:
		if value, ok := param.(float32); ok {
		data = decimal.NewFromFloat32(value)
	}else{
		err = errors.New("float32类型转换decimal失败")
	}
	default:
		err = errors.New("类型不支持计算")
	}
	return data, err
}