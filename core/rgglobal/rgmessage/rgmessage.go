package rgmessage

import (
	"github.com/jackylee92/rgo/core/rgconfig"
	"github.com/jackylee92/rgo/core/rgglobal/rgconst"
	"github.com/jackylee92/rgo/core/rgglobal/rgerror"
)

var rgMsg map[int64]map[string]string = map[int64]map[string]string{
	rgconst.ReturnSuccessCode: {
		"zh": rgerror.CurlSuccess,
		"en": rgerror.CurlSuccessEn,
	},
	rgconst.ReturnPanicCode: {
		"zh": rgerror.CurlErrorServerSelf,
		"en": rgerror.CurlErrorServerSelfEn,
	},
}

/*
* @Content : 初始化语言包
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2022-03-01
 */
func Msg(code int64) string {
	language := rgconfig.GetStr(rgconst.ConfigKeyMessage)
	if language != "zh" && language != "en" {
		language = "zh"
	}
	msgArr, ok := rgMsg[code]
	msg, ok := msgArr[language]
	if !ok {
		return rgconst.UnknownError
	}
	return msg
}

/*
* @Content : 注入项目语言包
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2022-03-10
 */
func InitAppMsg(param map[int64]string) {
	for code, msg := range param {
		rgMsg[code] = map[string]string{
			"en": msg,
			"zh": msg,
		}
	}
}
