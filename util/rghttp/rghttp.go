package rghttp

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"rgo"
	"strconv"
	"strings"
	"time"

	"rgo/core/rgconfig"
	"rgo/core/rgglobal/rgconst"
	"rgo/core/rgrequest"
)

type Client struct {
	Param    interface{}
	Method   string
	Header   map[string]string
	Url      string
	This     *rgrequest.Client `json:"-"` // json解析忽略
	duration int64             // 请求耗时 毫秒
}

const (
	headerUniqIDKey = rgconst.ContextUniqIDKey
	configCurlLog   = "util_curl_log" // curl日志是否开启记录
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
func (c *Client) getClient() (client *http.Client, err error) {
	client = &http.Client{
		Timeout:   time.Second * 10, // 超时时间
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
func (c *Client) getClientHeader() (req *http.Request, err error) {
	if c.Method == "POST" {
		switch c.Param.(type) {
		case []byte:
			if value, ok := c.Param.([]byte); ok {
				param := bytes.NewBuffer(value)
				req, err = http.NewRequest(c.Method, c.Url, param)
			} else {
				return req, errors.New("请求参数异常")
			}
		case string:
			newValue1, ok := c.Param.(string)
			if !ok {
				return req, errors.New("请求参数异常")
			}
			value := []byte(newValue1)
			param := bytes.NewBuffer(value)
			req, err = http.NewRequest(c.Method, c.Url, param)
			// <LiJunDong : 2021-10-09 11:23:51> --- 其他form格式可以断言后再添加
		case url.Values:
			newParam, ok := c.Param.(url.Values)
			if !ok {
				return req, errors.New("请求表单异常")
			}
			param := strings.NewReader(newParam.Encode())

			req, err = http.NewRequest(c.Method, c.Url, param)
		case map[string]string:
			newParam, ok := c.Param.(map[string]string)
			if !ok {
				return req, errors.New("请求数组异常")
			}
			param := new(bytes.Buffer)
			w := multipart.NewWriter(param)
			for k, v := range newParam {
				w.WriteField(k, v)
			}
			w.Close()
			req, err = http.NewRequest(c.Method, c.Url, param)
			req.Header.Set("Content-Type", w.FormDataContentType())

		default:
			return req, errors.New("参数类型错误")
		}
	} else if c.Method == "GET" {
		req, err = http.NewRequest(c.Method, c.Url, nil)
	} else if c.Method == "DELETE" {
		req, err = http.NewRequest(c.Method, c.Url, nil)
	} else {
		return req, errors.New("Method错误|" + c.Method)
	}
	if err != nil {
		return req, errors.New("请求参数设置异常|" + err.Error())
	}
	if c.This != nil {
		req.Header.Set(headerUniqIDKey, c.This.UniqId)
	}else{
		req.Header.Set(headerUniqIDKey, bootstrap.This.UniqId)
	}
	if len(c.Header) != 0 {
		for headerTitle, headerValue := range c.Header {
			req.Header.Set(headerTitle, headerValue)
		}
	}
	return req, nil
}

/*
 * @Content : http请求
 * @Param   : nil
 * @Return  : string
 * @Author  : LiJunDong
 * @Time    : 2020/7/22
 */
func (c *Client) GetApi() (data string, err error) {
	if c.Url == "" {
		return "", errors.New("URL不能为空")
	}
	if c.Method == "" {
		return "", errors.New("Method不能为空")
	}
	maxReqCh <- struct{}{}
	defer func() {
		<-maxReqCh
	}()
	client, err := c.getClient()
	if err != nil {
		return "", errors.New("获取请求客户端失败|" + err.Error())
	}
	//获取请求头
	req, err := c.getClientHeader()
	if err != nil {
		return "", errors.New("获取请求头失败|" + err.Error())
	}
	if rgconfig.GetBool(configCurlLog) {
		logStr, _ := json.Marshal(c)
		if c.This != nil {
			c.This.Log.Info("HTTP请求开始|" + string(logStr))
		}else{
			bootstrap.This.Log.Info("HTTP请求开始|" + string(logStr))
		}
	}
	startTimeInt := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		return "", errors.New("请求失败|" + err.Error())
	}
	endTimeInt := time.Now()
	processTime := (endTimeInt.UnixNano() - startTimeInt.UnixNano()) / 1000000 // 纳秒转毫秒
	c.duration = processTime
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("请求返回数据处理失败|" + err.Error())
	}
	if rgconfig.GetBool(configCurlLog) {
		logStr, _ := json.Marshal(map[string]interface{}{"url": c.Url, "httpStatus": resp.StatusCode, "body": string(body), "duration": c.duration})
		if c.This != nil {
			c.This.Log.Info("HTTP请求结束|" + string(logStr))
		}else{
			bootstrap.This.Log.Info("HTTP请求结束|" + string(logStr))
		}
	}
	if resp.StatusCode != 200 {
		return "", errors.New("请求相应HttpCode异常" + strconv.Itoa(resp.StatusCode))
	}
	return string(body), nil
}
