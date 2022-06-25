package rgrequest

import (
	"bytes"
	"go.opentelemetry.io/otel/bridge/opentracing"
	"io/ioutil"
	"rgo/core/rgglobal/rgconst"
	"time"

	"rgo/core/rgjaerger"
	"rgo/core/rglog"
	"rgo/core/rgmodel/rgmysql"
	"rgo/core/rgresponse"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

type Client struct {
	Param          interface{}
	Log            *rglog.Client
	UniqId         string
	Response       *rgresponse.Client
	requestTimeInt int64
	requestEndTime int64
	Jaerger        *rgjaerger.Client
	Mysql          *rgmysql.Factory
	Ctx            *gin.Context
}

/*
* @Content : 获取单次请求对象
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2022-03-11
 */
func Get(c *gin.Context) *Client {
	uniqId := GetUniqId(c)
	this, ok := c.Get(rgconst.ContextContainerKey)
	if ok {
		return this.(*Client)
	}
	ctxInterface, ok := c.Get(rgconst.ContextJeargerCtxKey)
	tracerInterface, _ := c.Get(rgconst.ContextJeargerTracerKey)
	jeargerClient := &rgjaerger.Client{}
	if ok {
		parentCtx := ctxInterface.(opentracing.SpanContext)
		tracer := tracerInterface.(opentracing.Tracer)
		jeargerClient = rgjaerger.New(parentCtx, tracer)
	}
	logger := rglog.New(uniqId)
	c.Set(rgconst.ContextStartTimeKey, time.Now().UnixNano())
	return &Client{
		Log:      logger,
		UniqId:   uniqId,
		Response: rgresponse.New(uniqId, c),
		Jaerger:  jeargerClient,
		Mysql:    rgmysql.New(uniqId, logger),
		Ctx:      c,
	}
}

/*
* @Content : 获取post json
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2022-03-28
 */
func GetPostJson(c *gin.Context) (data string) {
	dataInterface, ok := c.Get("post_json")
	if ok {
		return dataInterface.(string)
	}
	var bodyBytes []byte
	// 从原有Request.Body读取
	bodyBytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return data
	}
	// 新建缓冲区并替换原有Request.body
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	data = string(bodyBytes)
	c.Set("post_json", data)
	return data

}

/*
* @Content : 获取uniqid
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2022-03-29
 */
func GetUniqId(c *gin.Context) string {
	value, ok := c.Get(rgconst.ContextUniqIDKey)
	if !ok {
		var uniqid string
		if c.Request != nil {
			uniqid = c.GetHeader(rgconst.ContextUniqIDKey)
		}
		if uniqid == "" {
			uniqid = xid.New().String()
		}
		c.Set(rgconst.ContextUniqIDKey, uniqid)
		return uniqid
	} else {
		return value.(string)
	}
}
