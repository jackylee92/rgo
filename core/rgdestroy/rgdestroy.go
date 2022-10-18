package rgdestroy

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/jackylee92/rgo/core/rgglobal/rgconst"
	"github.com/jackylee92/rgo/core/rglog"
)

func Listen() {
	//  用于系统信号的监听
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM) // 监听可能的退出信号
		received := <-c                                                                           //接收信号管道中的值
		rglog.SystemInfo(rgconst.ProcessKilled, "信号值"+":"+received.String())
		close(c)
		os.Exit(1)
	}()

}
