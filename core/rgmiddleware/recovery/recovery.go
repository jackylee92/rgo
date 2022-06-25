package recovery

import (
	"errors"

	"rgo/core/rgglobal/rgerror"
	"rgo/core/rglog"
	"rgo/core/rgresponse"

	"github.com/gin-gonic/gin"
)

type PanicExceptionRecord struct{}

// CustomRecovery 自定义错误(panic等)拦截中间件、对可能发生的错误进行拦截、统一记录
func Handle() gin.HandlerFunc {
	DefaultErrorWriter := &PanicExceptionRecord{}
	return gin.RecoveryWithWriter(DefaultErrorWriter, func(c *gin.Context, err interface{}) {
		// 这里针对发生的panic等异常进行统一响应即可
		// 这里的 err 数据类型为 ：runtime.boundsError  ，需要转为普通数据类型才可以输出
		rgresponse.SystemError(c, nil)
	})
}

func (p *PanicExceptionRecord) Write(b []byte) (n int, err error) {
	errStr := string(b)
	err = errors.New(errStr)
	rglog.SystemError(rgerror.ErrorRequirePanic + "|err:" + errStr)
	return len(errStr), err
}
