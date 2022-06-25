package requestlog

import (
	"net/http"
	"net/url"

	"rgo/core/rgconfig"
	"rgo/core/rgglobal/rgconst"
	"rgo/core/rgjson"
	"rgo/core/rglog"
	"rgo/core/rgrequest"

	"github.com/gin-gonic/gin"
)

type LogS struct {
	IP       string      `json:"ip"`
	ClientIp string      `json:"client_ip"`
	Url      string      `json:"url"`
	Headers  http.Header `json:"header"`
	FormPost url.Values  `json:"form_post"`
	FormJson string      `json:"form_json"`
	Time     string      `json:"time"`
	Method   string      `json:"method"`
}

/*
* @Content : 请求日志
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2022-03-10
 */
func Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		logLeve := rgconfig.GetStr(rgconst.ConfigKeyRequestLog)
		if logLeve == "1" || logLeve == "2" {
			requestData := LogS{
				ClientIp: c.ClientIP(),
				IP:       c.Request.RemoteAddr,
				Url:      c.Request.RequestURI,
				Headers:  c.Request.Header,
				FormPost: c.Request.PostForm,
				FormJson: GetPostJson(c),
				Method:   c.Request.Method,
			}
			logData, _ := rgjson.Marshel(requestData)
			uniqId := rgrequest.GetUniqId(c)
			rglog.RequestLog(uniqId, "REQUEST", logData)
		}
		c.Next()
	}

}

/*
* @Content : 获取Post json参数
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2021-09-10
 */
func GetPostJson(c *gin.Context) (data string) {
	return rgrequest.GetPostJson(c)
}
