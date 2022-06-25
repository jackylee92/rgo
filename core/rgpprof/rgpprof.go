package rgpprof

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"rgo/core/rgconfig"
	"rgo/core/rgglobal/rgconst"
	"rgo/core/rglog"
)

/*
* @Content : 启动压测调试工具pprof
* 在有压力期间 go tool pprof http://端口:ip/debug/pprof/profile
* cd ~/pprof
* go tool pprof -http 0.0.0.0:3001 指定当前目录下收集的文件
*
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2022-05-20
 */
func Start(router *gin.Engine) {
	if rgconfig.GetBool(rgconst.ConfigPProf) {
		pprof.Register(router)
		rglog.SystemInfo("启动项【pprof】:成功")
	}
}
