package container

import (
	"rgo/core/rgglobal/rgconst"
	"rgo/core/rgrequest"

	"github.com/gin-gonic/gin"
)

/*
* @Content : 请求日志
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2022-03-10
 */
func Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(rgconst.ContextContainerKey, rgrequest.Get(c))
		c.Next()
	}

}
