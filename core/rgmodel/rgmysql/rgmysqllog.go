package rgmysql

import (
	"context"
	"errors"
	"strconv"
	"time"

	"rgo/core/rglog"

	"gorm.io/gorm/logger"
)

type Logger struct {
	logger.Writer
	logger.Config
	infoStr, warnStr, errStr            string
	traceStr, traceErrStr, traceWarnStr string
}

var ErrRecordNotFound = errors.New("record not found")

// LogMode log mode
func (l Logger) LogMode(level logger.LogLevel) logger.Interface {
	l.LogLevel = level
	return l
}

// Info print info
func (l Logger) Info(ctx context.Context, msg string, data ...interface{}) {
	if ctx.Value("logger") == nil {
		return
	}
	loggerClient := ctx.Value("logger").(*rglog.Client)
	loggerClient.Info(msg)
}

// Warn print warn messages
func (l Logger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if ctx.Value("logger") == nil {
		return
	}
	loggerClient := ctx.Value("logger").(*rglog.Client)
	loggerClient.Warn(msg)
}

// Error print error messages
func (l Logger) Error(ctx context.Context, msg string, data ...interface{}) {
	if ctx.Value("logger") == nil {
		return
	}
	loggerClient := ctx.Value("logger").(*rglog.Client)
	loggerClient.Error(msg)
}

// Trace print sql message
func (l Logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()
	// 出错误
	if err != nil {
		l.Error(ctx, sql+" "+err.Error())
		return
	}
	if elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= logger.Warn {
		l.Warn(ctx, "慢查询:"+sql+" 影响:"+strconv.Itoa(int(rows))+" 耗时:"+elapsed.String())
	}
	if l.LogLevel >= logger.Info {
		l.Info(ctx, sql+" 影响:"+strconv.Itoa(int(rows))+" 耗时:"+elapsed.String())
	}
	return
}
