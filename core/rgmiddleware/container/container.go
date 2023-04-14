package container

import (
	"github.com/jackylee92/rgo/core/rgglobal/rgconst"
	"github.com/jackylee92/rgo/core/rgrequest"

	"github.com/gin-gonic/gin"
)

func Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(rgconst.ContextContainerKey, rgrequest.Get(c))
		c.Next()
	}

}
