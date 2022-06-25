package rgstarthook

import (
	"rgo/core/rgglobal/rgconst"
	"rgo/core/rgglobal/rgerror"
	"rgo/core/rglog"
	"sync"
)

/*
 * @Content : rgstart
 * @Author  : LiJunDong
 * @Time    : 2022-05-30$
 */


// startFuncList <LiJunDong : 2022-05-30 18:01:49> --- 启动完成后执行函数列表
var startFuncList = make([]func(), 0, 20)
var lock sync.Mutex



// RegisterStartFunc 注册启动完成后执行的函数
// @Param   : f func()
// @Return  :
// @Author  : LiJunDong
// @Time    : 2022-05-30
func RegisterStartFunc(f func())  {
	lock.Lock()
	if len(startFuncList) >= rgconst.StartFuncMaxCount {
		rglog.SystemError(rgerror.ErrorStartFuncOverflow)
		return
	}
	startFuncList = append(startFuncList, f)
	lock.Unlock()
}

// Run 服务启动完成后，执行用户注册的函数
// @Param   :
// @Return  :
// @Author  : LiJunDong
// @Time    : 2022-05-30
func Run()  {
	for _,f := range startFuncList {
 		f()
	}
}
