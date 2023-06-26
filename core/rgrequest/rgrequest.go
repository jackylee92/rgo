package rgrequest

import (
	"bytes"
	"github.com/jackylee92/rgo/core/rgglobal/rgconst"
	"github.com/jackylee92/rgo/core/rgjaeger"
	"github.com/jackylee92/rgo/core/rgmysql"
	"io"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"

	"github.com/jackylee92/rgo/core/rglog"
	"github.com/jackylee92/rgo/core/rgresponse"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/rs/xid"
)

type Client struct {
	Param          interface{}
	Log            *rglog.Client
	UniqId         string
	Response       *rgresponse.Client
	requestTimeInt int64
	requestEndTime int64
	Jaeger         *rgjaeger.Client
	Mysql          *rgmysql.Factory
	Ctx            *gin.Context
}

var ClientPool sync.Pool

func init() {
	ClientPool.New = func() any {
		nilUniqId := ""
		logClient := rglog.New(nilUniqId)
		return &Client{
			Log:      logClient,
			UniqId:   nilUniqId,
			Mysql:    rgmysql.New(nilUniqId, logClient),
			Response: &rgresponse.Client{},
			Jaeger:   &rgjaeger.Client{},
		}
	}
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
	ctxInterface, ok := c.Get(rgconst.ContextJaegerCtxKey)
	tracerInterface, _ := c.Get(rgconst.ContextJaegerTracerKey)
	jaegerClient := &rgjaeger.Client{}
	if ok {
		parentCtx := ctxInterface.(opentracing.SpanContext)
		tracer := tracerInterface.(opentracing.Tracer)
		jaegerClient.Reset(parentCtx, tracer)
	}
	c.Set(rgconst.ContextStartTimeKey, time.Now().UnixNano())
	client := ClientPool.Get().(*Client)
	client.UniqId = uniqId
	client.Ctx = c
	client.Jaeger = jaegerClient
	client.reset()
	return client
}

func (c *Client) reset() {
	c.Log.UniqId = c.UniqId
	c.Mysql.UniqId = c.UniqId
	c.Response.UniqId = c.UniqId
	c.Response.Ctx = c.Ctx
}

func GetPostJson(c *gin.Context) (data string) {
	dataInterface, ok := c.Get("post_json")
	if ok {
		return dataInterface.(string)
	}
	var bodyBytes []byte
	// 从原有Request.Body读取
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return data
	}
	// 新建缓冲区并替换原有Request.body
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
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

// Proxy <LiJunDong : 2023-01-12 15:14:07> --- 请求代理转发
func (c *Client) Proxy(scheme, host, path string) (err error) {
	var target = scheme + "://" + host + path
	proxyUrl, err := url.Parse(target)
	if err != nil {
		c.Log.Error("代理失败", target, err)
		return err
	}
	proxyUrl.Scheme = scheme
	proxyUrl.Host = host
	c.Ctx.Request.Host = host
	c.Ctx.Request.URL.Path = path
	proxy := httputil.NewSingleHostReverseProxy(proxyUrl)
	proxy.ServeHTTP(c.Ctx.Writer, c.Ctx.Request)
	c.Log.Info("代理结束", target)
	return err
}
