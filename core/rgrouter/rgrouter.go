package rgrouter

import (
	"github.com/gin-gonic/gin"
	"github.com/jackylee92/rgo/core/rgconfig"
	"github.com/jackylee92/rgo/core/rgglobal/rgconst"
	"github.com/jackylee92/rgo/core/rgglobal/rgerror"
	"github.com/jackylee92/rgo/core/rglog"
	"github.com/jackylee92/rgo/core/rgmiddleware/container"
	"github.com/jackylee92/rgo/core/rgmiddleware/crossdomain"
	"github.com/jackylee92/rgo/core/rgmiddleware/jeager"
	"github.com/jackylee92/rgo/core/rgmiddleware/recovery"
	"github.com/jackylee92/rgo/core/rgmiddleware/requestlog"
	"github.com/jackylee92/rgo/core/rgpprof"
	"io"
	"net/http"
)

/*
* @Content : 初始化路由
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2022-03-01
 */
func NewRouter() *gin.Engine {
	_ = InitTrans()
	router := &gin.Engine{}
	if rgconfig.GetBool(rgconst.ConfigKeyDebug) {
		gin.SetMode(gin.ReleaseMode) // <LiJunDong : 2022-06-02 18:26:14> --- 关闭gin的很长一段提示信息
		router = gin.New()
		rglog.SystemInfo("启动项【debug】:成功")
	} else {
		gin.SetMode(gin.ReleaseMode)
		gin.DefaultWriter = io.Discard
		router = gin.New()
	}
	if rgconfig.GetBool(rgconst.ConfigKeyHttpAllowCrossDomain) {
		router.Use(crossdomain.Handle()) // 跨域
	}
	router.Use(jeager.Handle(), requestlog.Handle(), recovery.Handle(), container.Handle())
	router.GET("/"+rgconst.ConfigHeartBeatUrl, HeartBeatHandle)     // 健康检查
	router.GET("/"+rgconst.ConfigSetLogLevelUrl, SetLogLevelHandle) // 设置日志级别
	router.GET("/"+rgconst.ConfigGetLogLevelUrl, GetLogLevelHandle) // 获取日志级别
	router.GET("/"+rgconst.ConfigGetUrl, GetConfigHandle)           // 获取日志配置

	router.NoRoute(func(c *gin.Context) { c.String(http.StatusNotFound, rgerror.Curl404Error) })

	return router
}

/*
* @Content : 启动
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2022-05-21
 */
func Run(router *gin.Engine) {
	rgpprof.Start(router)
	rglog.SystemInfo("启动项【port】:" + rgconfig.GetStr(rgconst.ConfigKeyPort))
	router.Run(":" + rgconfig.GetStr(rgconst.ConfigKeyPort))
}
