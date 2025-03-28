package rgo

import (
	"github.com/gin-gonic/gin"
	"github.com/jackylee92/rgo/core/rgconfig"
	"github.com/jackylee92/rgo/core/rgconfig/rgapollo"
	"github.com/jackylee92/rgo/core/rgconfig/rglocal"
	"github.com/jackylee92/rgo/core/rgdestroy"
	"github.com/jackylee92/rgo/core/rgglobal"
	"github.com/jackylee92/rgo/core/rgglobal/rgconst"
	"github.com/jackylee92/rgo/core/rgjaeger"
	"github.com/jackylee92/rgo/core/rglog"
	"github.com/jackylee92/rgo/core/rgmysql"
	"github.com/jackylee92/rgo/core/rgrequest"
	"github.com/jackylee92/rgo/core/rgrequired"
	"github.com/jackylee92/rgo/core/rgstarthook"
	"log"
)

/*
* @Content : 项目启动加载核型集合
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2022-03-01
 */
func init() {
	// <LiJunDong : 2022-03-11 11:09:35> --- 加载系统常量配置
	rgglobal.Load()
	// <LiJunDong : 2022-03-11 11:09:35> --- 注入本地配置
	rglocal.Register()
	// <LiJunDong : 2022-03-11 11:09:35> --- 注入Apollo配置
	rgapollo.Register()
	// <LiJunDong : 2022-03-11 11:09:35> --- 启动配置
	rgconfig.Start()
	// <LiJunDong : 2022-03-11 11:09:35> --- 监听推出信号
	rgdestroy.Listen()
	// <LiJunDong : 2022-03-11 11:09:35> --- 开启日志
	rglog.Start()
	// <LiJunDong : 2022-03-11 11:09:35> --- 项目依赖检查
	rgrequired.Check()
	// <LiJunDong : 2022-03-21 11:09:35> --- 设置项目名称，配置中必须包含
	rgglobal.SetAppName(rgconfig.GetStr(rgconst.ConfigKeyAppName))
	// <LiJunDong : 2022-03-21 11:09:35> --- 设置jaeger是否开启
	rgjaeger.SetJaegerStatus(rgconfig.GetBool(rgconst.ConfigKeyJaegerStatus))
	// <LiJunDong : 2022-03-30 21:18:37> --- 根据配置是否启用mysql
	rgmysql.Start()
	// <LiJunDong : 2022-05-30 18:14:09> --- 启动用户注册的函数
	rgstarthook.Run()
	//  <LiJunDong : 2022-06-21 21:21:04> --- 启动一个服务自身的容器
	serverContainer()
	startMsg := "启动项【rgo-init】:成功"
	rglog.SystemInfo(startMsg)
	log.Println("|SystemInfo| " + startMsg + "| UniqId:终端")
}

var This *rgrequest.Client

// serverContainer 服务启动自身的一个对象，在非web请求，无上下文this对象时使用这个
// @Param   :
// @Return  : this *rgrequest.Client
// @Author  : LiJunDong
// @Time    : 2022-06-21
func serverContainer() {
	ctx := &gin.Context{}
	ctx.Set(rgconst.ContextUniqIDKey, "Init")
	This = rgrequest.Get(ctx)
}
