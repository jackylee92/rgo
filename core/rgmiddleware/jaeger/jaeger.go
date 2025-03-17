package jaeger

import (
	"bytes"
	"encoding/json"
	"github.com/jackylee92/rgo/core/rgglobal/rgconst"
	"io"
	"strconv"
	"time"

	"github.com/jackylee92/rgo/core/rgjaeger"
	"github.com/jackylee92/rgo/core/rglog"
	"github.com/jackylee92/rgo/core/rgrequest"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go/ext"
)

// bodyLogWriter是为了记录返回数据到log中进行了双写
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// 继承ResponseWriter，重写write方法
func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		if rgjaeger.JaegerStatus() {
			tracer, spanContext, closer, err := rgjaeger.GetTracer(c.Request.Header)
			if err == nil {
				defer func(closer io.Closer) {
					_ = closer.Close()
				}(closer)
				startSpan := tracer.StartSpan(c.Request.URL.Path, ext.RPCServerOption(spanContext))
				defer startSpan.Finish()
				blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
				c.Writer = blw // 注入自己重写ResponseWriter的对象

				var tagName ext.StringTagName
				startTimeInt := time.Now()
				startTime := startTimeInt.Format(rgconst.GoTimeFormat)

				ext.HTTPUrl.Set(startSpan, c.Request.URL.String())
				ext.HTTPMethod.Set(startSpan, c.Request.Method)
				// header
				header := getRequestHeader(c)
				if header != "" {
					tagName = "http.header"
					tagName.Set(startSpan, header)
				}
				// post param
				postParam := getRequestParam(c)
				if postParam != "" {
					tagName = "http.request"
					tagName.Set(startSpan, postParam)
				}
				tagName = "time.start"
				tagName.Set(startSpan, startTime)
				parentCtx := startSpan.Context()
				c.Set(rgconst.ContextJaegerTracerKey, tracer)
				c.Set(rgconst.ContextJaegerCtxKey, parentCtx)
				c.Next()
				endTimeInt := time.Now()
				endTime := endTimeInt.Format(rgconst.GoTimeFormat)
				tagName = "time.end"
				tagName.Set(startSpan, endTime)
				tagName = "time.duration"
				processTime := (endTimeInt.UnixNano() - startTimeInt.UnixNano()) / 1000000 // 毫秒
				tagName.Set(startSpan, strconv.Itoa(int(processTime)))
				ext.HTTPStatusCode.Set(startSpan, uint16(c.Writer.Status()))
				resp := blw.body.String()
				tagName = "http.response"
				tagName.Set(startSpan, resp)
			} else {
				rglog.SystemError(err.Error())
				c.Next()
			}
		} else {
			c.Next()
		}
	}
}

/*
* @Content : 获取所有post参数
* @Param   :
* @Return  :
* @Author  : `!v g:USER`
* @Time    : `!v strftime('%Y-%m-%d')`
 */
func getRequestParam(c *gin.Context) string {
	var paramForm, paramJson string
	_ = c.Request.ParseForm()
	paramFormArr := make(map[string]interface{}, len(c.Request.PostForm)+5)
	for k, v := range c.Request.PostForm {
		paramFormArr[k] = v
	}
	if len(paramFormArr) != 0 {
		paramFormByte, _ := json.Marshal(paramFormArr)
		paramForm = string(paramFormByte)
	}
	paramJson = rgrequest.GetPostJson(c)
	if len(paramFormArr) == 0 {
		return paramJson
	}
	if len(paramFormArr) != 0 && paramJson == "" {
		return paramForm
	}
	if len(paramFormArr) != 0 && paramJson != "" {
		return "form:" + paramForm + " | json:" + paramJson
	}
	return "form:" + paramForm + " | json:" + paramJson
}

/*
* @Content : 获取所有header
* @Param   :
* @Return  :
* @Author  : `!v g:USER`
* @Time    : `!v strftime('%Y-%m-%d')`
 */
func getRequestHeader(c *gin.Context) string {
	header := c.Request.Header
	headerJson, err := json.Marshal(header)
	if err != nil {
		return ""
	}
	return string(headerJson)
}
