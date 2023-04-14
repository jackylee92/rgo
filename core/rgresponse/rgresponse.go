package rgresponse

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/jackylee92/rgo/core/rgconfig"
	"github.com/jackylee92/rgo/core/rgglobal/rgconst"
	"github.com/jackylee92/rgo/core/rgglobal/rgmessage"
	"github.com/jackylee92/rgo/core/rglog"
	"net/http"
	"time"
)

type Client struct {
	UniqId string
	Ctx    *gin.Context
}

// <LiJunDong : 2022-03-10 18:59:25> --- 所有返回相关函数

func returnJson(Context *gin.Context, dataCode int64, data, msg interface{}) {
	if Context == nil || Context.Request == nil {
		return
	}
	if msg == nil {
		msg = rgmessage.Msg(dataCode)
	}
	returnData := map[string]interface{}{
		"code":    dataCode,
		"message": msg,
		"data":    data,
	}
	if rgconfig.GetBool(rgconst.ConfigReturnUniqId) {
		returnData["uniqId"] = Context.GetString(rgconst.ContextUniqIDKey)
	}

	if rgconfig.GetStr(rgconst.ConfigKeyRequestLog) == "2" {
		startTime := Context.GetInt64(rgconst.ContextStartTimeKey)
		logData := returnData
		if startTime != 0 {
			logData["duration"] = (time.Now().UnixNano() - startTime) / 1000000
		}
		logStr, _ := json.Marshal(logData)
		rglog.RequestLog(Context.GetString(rgconst.ContextUniqIDKey), "RESPONSE", string(logStr))
		delete(returnData, "duration")
	}
	var ginReturnData gin.H
	ginReturnData = returnData
	//Context.Header("key2020","value2020")     //可以根据实际情况在头部添加额外的其他信息
	Context.Header(rgconst.ContextUniqIDKey, Context.GetString(rgconst.ContextUniqIDKey)) //可以根据实际情况在头部添加额外的其他信息
	Context.JSON(http.StatusOK, ginReturnData)
	Context.Abort()
}

// <LiJunDong : 2022-03-10 18:59:25> --- 所有返回相关函数

// 语法糖函数封装

// 系统异常
func SystemError(c *gin.Context, data interface{}) {
	returnJson(c, http.StatusInternalServerError, data, nil)
}

/*
* @Content : 正确返回，固定返回code为200
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2022-03-13
 */
func (c *Client) ReturnSuccess(data interface{}) {
	returnJson(c.Ctx, rgconst.ReturnSuccessCode, data, nil)
}

/*
* @Content : 自定义返回 code, data, message
* @Param   : code, data, message
* @Return  :
* @Author  : LiJunDong
* @Time    : 2022-03-13
 */
func (c *Client) Return(errCode int64, returnData ...interface{}) {
	var data, msg interface{}
	if len(returnData) > 0 {
		data = returnData[0]
	}
	if len(returnData) > 1 {
		msg = returnData[1]
	}
	returnJson(c.Ctx, errCode, data, msg)
	return
}

/*
* @Content : 错误返回
* @Param   : code, data, message
* @Return  :
* @Author  : LiJunDong
* @Time    : 2022-03-13
 */
func (c *Client) ReturnError(errCode int64, returnData ...interface{}) {
	var data, msg interface{}
	if len(returnData) > 0 {
		data = returnData[0]
	}
	if len(returnData) > 1 {
		msg = returnData[1]
	}
	returnJson(c.Ctx, errCode, data, msg)
	return
}

func (c *Client) View(path string, data map[string]interface{}) {
	if data == nil {
		data = gin.H{}
	}
	c.Ctx.HTML(http.StatusOK, path, data)
	c.Ctx.Abort()
}

func New(uniqId string, ctx *gin.Context) *Client {
	return &Client{
		UniqId: uniqId,
		Ctx:    ctx,
	}
}
