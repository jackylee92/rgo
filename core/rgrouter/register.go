package rgrouter

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/jackylee92/rgo/core/rgconfig"
	"github.com/jackylee92/rgo/core/rgglobal/rgconst"
	"github.com/jackylee92/rgo/core/rglog"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

var registerFuncMap map[string]func(validator.FieldLevel) bool
var registerMessageMap map[string]MessageTrans
var Validate = validator.New()
var Trans ut.Translator

type MessageTrans struct {
	RegisterFn    func(ut.Translator) error
	TranslationFn func(ut.Translator, validator.FieldError) string
}

/*
* @Content : init
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2022-03-29
 */
func init() {
	registerFuncMap = make(map[string]func(validator.FieldLevel) bool, 0)
	registerMessageMap = make(map[string]MessageTrans, 0)
}

// <LiJunDong : 2022-03-29 17:57:53> --- 注册语言包
// <LiJunDong : 2022-03-29 17:57:59> --- 注册自定义验证

func InitTrans() (err error) {
	var ok bool
	if Validate, ok = binding.Validator.Engine().(*validator.Validate); ok {
		for key, item := range registerFuncMap {
			_ = Validate.RegisterValidation(key, item)
		}
		zhT := zh.New() //中文翻译器
		enT := en.New() //英文翻译器
		uni := ut.New(enT, zhT, enT)

		locale := rgconfig.GetStr(rgconst.ConfigKeyMessage)
		// 也可以使用 uni.FindTranslator(...) 传入多个locale进行查找
		Trans, ok = uni.GetTranslator(locale)
		if !ok {
			rglog.SystemError("注册自定义验证语言包失败")
			return
		}

		for key, item := range registerMessageMap {
			_ = Validate.RegisterTranslation(key, Trans, item.RegisterFn, item.TranslationFn)
		}

		Validate.RegisterTagNameFunc(func(field reflect.StructField) string {
			label := field.Tag.Get("label")
			if label == "" {
				label = field.Tag.Get("json")
			}
			if label == "" {
				label = field.Tag.Get("form")
			}
			if label == "" {
				return field.Name
			}
			return label
		})
		// 注册翻译器
		switch locale {
		case "en":
			err = en_translations.RegisterDefaultTranslations(Validate, Trans)
		case "zh":
			err = zh_translations.RegisterDefaultTranslations(Validate, Trans)
		default:
			err = en_translations.RegisterDefaultTranslations(Validate, Trans)
		}
		return
	} else {
		rglog.SystemError("注册自定义验证失败")
	}
	return
}

/*
* @Content : 自定义验证方法
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2022-03-29
 */
func RegisterFunc(name string, f func(validator.FieldLevel) bool) {
	registerFuncMap[name] = f
}

/*
* @Content : 注册自定义错误语言
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2022-03-29
 */
func RegisterMessage(name string, param MessageTrans) {
	registerMessageMap[name] = param
}

/*
* @Content : 返回入参错误
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2022-04-15
 */
func Error(err error) (msg string, fields []string) {
	jsonErr, ok := err.(*json.UnmarshalTypeError)
	if ok {
		if jsonErr.Struct != "" || jsonErr.Field != "" {
			msg = "参数[" + jsonErr.Field + "]错误，要求[" + jsonErr.Type.String() + "]类型"
			// "json: cannot unmarshal " + e.Value + " into Go struct field " + e.Struct + "." + e.Field + " of type " + e.Type.String()
			// json: cannot unmarshal string into Go struct field DetailUserParam.id of type int64
		} else {
			msg = "请求参数解析错误"
			// "json: cannot unmarshal " + e.Value + " into Go value of type " + e.Type.String()
		}
	} else {
		validateErrs, ok := err.(validator.ValidationErrors)
		if ok {
			errList := make([]string, 0, len(validateErrs))
			for _, e := range validateErrs {
				errList = append(errList, e.Translate(Trans))
				fields = append(fields, e.Field())
			}
			msg = strings.Join(errList, ",")
		} else {
			msg = "请求参数格式错误"
		}
	}
	return msg, fields
}

// HeartBeatHandle 心跳
// @Param   : ctx *gin.Context
// @Return  :
// @Author  : LiJunDong
// @Time    : 2022-06-18
func HeartBeatHandle(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, nil)
}

func SetLogLevelHandle(ctx *gin.Context) {
	param := ctx.Query("level")
	rglog.SetLogLever(param)
	ctx.JSON(http.StatusOK, map[string]int{"code": 200})
}

func GetLogLevelHandle(ctx *gin.Context) {
	level := rglog.GetLogLevel()
	ctx.JSON(http.StatusOK, map[string]string{"level": level})
}

func GetConfigHandle(ctx *gin.Context) {
	key := ctx.Query("key")
	if key == "" {
		ctx.JSON(http.StatusOK, map[string]string{"error": "key is nil"})
		return
	}
	value := rgconfig.Get(key)
	ctx.JSON(http.StatusOK, map[string]interface{}{key: value})
	return
}
