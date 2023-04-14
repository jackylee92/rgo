package rghttp

import (
	"bytes"
	"errors"
	"github.com/jackylee92/rgo"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/jackylee92/rgo/core/rgconfig"
	"github.com/jackylee92/rgo/core/rgglobal/rgconst"
	"github.com/jackylee92/rgo/core/rgrequest"
)

type Client struct {
	Param    interface{}
	Method   string
	Header   map[string]string
	Url      string
	Timeout  int               // 超时时间，如果没天默认10秒
	This     *rgrequest.Client `json:"-"` // json解析忽略
	duration int64             // 请求耗时 毫秒
}

const (
	headerUniqIDKey = rgconst.ContextUniqIDKey
	configCurlLog   = "util_curl_log" // curl日志是否开启记录
	defaultTimeOut  = 10              // 默认超时时间
)

var globalTransport *http.Transport
var maxReqCh chan struct{}

/*
* @Content : init
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2022-03-31
 */
func init() {
	globalTransport = http.DefaultTransport.(*http.Transport).Clone()
	globalTransport.ResponseHeaderTimeout = time.Second * time.Duration(10) // 限制了读取头部的时间
	globalTransport.MaxIdleConns = 100                                      // 最大空闲连接数量
	globalTransport.MaxIdleConnsPerHost = 100                               // 每个host的设定的空闲连接数量
	globalTransport.MaxConnsPerHost = 100                                   // 最大空闲连接数量
	maxReqCh = make(chan struct{}, 20000)                                   // 并发请求数
}

/*
 * @Content : 获取httpClient
 * @Param   : nil
 * @Return  : httpClient
 * @Author  : LiJunDong
 * @Time    : i
 */
func (i *Client) getClient() (client *http.Client, err error) {
	client = &http.Client{
		Timeout:   time.Second * time.Duration(i.Timeout), // 超时时间
		Transport: globalTransport,
	}
	return client, nil
}

/*
 * @Content : 获取请求头
 *  POST参数 strings.NewReader
 *      strings.newReader(data) data为string格式
 *      data可有两种 :
 *          json : data := "{...}" string格式
 *          form : data := url.Values{} 数组
 *              data[key] = []string{value}
 *              data.Encode()
 * @Param   : nil
 * @Return  : 返回介绍
 * @Author  : LiJunDong
 * @Time    : 2020/7/22
 */
func (i *Client) getClientHeader() (req *http.Request, err error) {
	if i.Method == "POST" {
		switch i.Param.(type) {
		case []byte:
			if value, ok := i.Param.([]byte); ok {
				param := bytes.NewBuffer(value)
				req, err = http.NewRequest(i.Method, i.Url, param)
			} else {
				return req, errors.New("请求参数异常")
			}
		case string:
			newValue1, ok := i.Param.(string)
			if !ok {
				return req, errors.New("请求参数异常")
			}
			value := []byte(newValue1)
			param := bytes.NewBuffer(value)
			req, err = http.NewRequest(i.Method, i.Url, param)
			// <LiJunDong : 2021-10-09 11:23:51> --- 其他form格式可以断言后再添加
		case url.Values:
			newParam, ok := i.Param.(url.Values)
			if !ok {
				return req, errors.New("请求表单异常")
			}
			param := strings.NewReader(newParam.Encode())

			req, err = http.NewRequest(i.Method, i.Url, param)
		case map[string]string:
			newParam, ok := i.Param.(map[string]string)
			if !ok {
				return req, errors.New("请求数组异常")
			}
			param := new(bytes.Buffer)
			w := multipart.NewWriter(param)
			for k, v := range newParam {
				_ = w.WriteField(k, v)
			}
			_ = w.Close()
			req, err = http.NewRequest(i.Method, i.Url, param)
			req.Header.Set("Content-Type", w.FormDataContentType())

		default:
			return req, errors.New("参数类型错误")
		}
	} else if i.Method == "GET" {
		req, err = http.NewRequest(i.Method, i.Url, nil)
	} else if i.Method == "DELETE" {
		req, err = http.NewRequest(i.Method, i.Url, nil)
	} else {
		return req, errors.New("Method错误|" + i.Method)
	}
	if err != nil {
		return req, errors.New("请求参数设置异常|" + err.Error())
	}
	if i.This != nil {
		req.Header.Set(headerUniqIDKey, i.This.UniqId)
	} else {
		req.Header.Set(headerUniqIDKey, rgo.This.UniqId)
	}
	if len(i.Header) != 0 {
		for headerTitle, headerValue := range i.Header {
			req.Header.Set(headerTitle, headerValue)
		}
	}
	return req, nil
}

func (i *Client) GetApi() (data string, err error) {
	if i.Url == "" {
		return "", errors.New("URL不能为空")
	}
	if i.Method == "" {
		return "", errors.New("method不能为空")
	}
	maxReqCh <- struct{}{}
	defer func() {
		<-maxReqCh
	}()
	if i.This == nil {
		i.This = rgo.This
	}
	i.addUniqId()
	if i.Timeout == 0 {
		i.Timeout = defaultTimeOut
	}
	client, err := i.getClient()
	if err != nil {
		return "", errors.New("获取请求客户端失败|" + err.Error())
	}
	//获取请求头
	req, err := i.getClientHeader()
	if err != nil {
		return "", errors.New("获取请求头失败|" + err.Error())
	}
	if rgconfig.GetBool(configCurlLog) {
		defer func() {
			i.This.Log.Info("CURL", i.duration, i.Url, i.Method, i.Timeout, i.Param, "RESULT: "+data, err)
			if err != nil {
				i.This.Log.Error("请求失败", err)
			}
		}()
	}
	startTimeInt := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		return "", errors.New("请求失败|" + err.Error())
	}
	endTimeInt := time.Now()
	processTime := (endTimeInt.UnixNano() - startTimeInt.UnixNano()) / 1000000 // 纳秒转毫秒 1000毫秒=1秒
	i.duration = processTime
	defer func() {
		if err := resp.Body.Close(); err != nil {
			i.This.Log.Error("请求结束关闭链接失败", i.Url, err)
		}
	}()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("请求返回数据处理失败|" + err.Error())
	}
	if rgconfig.GetBool(configCurlLog) {
		logStr := map[string]interface{}{"url": i.Url, "httpStatus": resp.StatusCode, "body": string(body), "duration": i.duration}
		if i.This != nil {
			i.This.Log.Info("HTTP请求结束", logStr)
		} else {
			rgo.This.Log.Info("HTTP请求结束", logStr)
		}
	}
	if resp.StatusCode != 200 {
		return "", errors.New("请求相应HttpCode异常" + strconv.Itoa(resp.StatusCode))
	}
	return string(body), nil
}

/*
 * @Content : common
 * @Author  : LiJunDong
 * @Time    : 2023-01-10$
 */

// Proxy 转发
// @Param   : scheme: "http" host:"127.0.0.1" path:"/api/get"
// @Return  : nil
// @Author  : LiJunDong
// @Time    : 2023-02-09
func (i *Client) Proxy(scheme, host, path string) {
	target := scheme + "://" + host
	proxyUrl, err := url.Parse(target)
	if err != nil {
		i.This.Log.Error("代理失败", target)
		return
	}
	proxyUrl.Scheme = scheme
	proxyUrl.Host = host
	i.This.Ctx.Request.Host = host
	i.This.Ctx.Request.URL.Path = path
	proxy := httputil.NewSingleHostReverseProxy(proxyUrl)
	proxy.ServeHTTP(i.This.Ctx.Writer, i.This.Ctx.Request)
	i.This.Log.Info("proxy", "代理结束", target)
}

func (i *Client) addUniqId() {
	hasParam := strings.Index(i.Url, "?")
	if hasParam == -1 {
		i.Url += "?"
	} else {
		i.Url += "&"
	}
	i.Url += headerUniqIDKey + "=" + i.This.UniqId
	return
}
