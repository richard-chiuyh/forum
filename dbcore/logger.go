package dbcore

import (
	"context"
	"errors"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

const slowThreshold = 200 * time.Millisecond

type Logger struct {
}

var _ logger.Interface = (*Logger)(nil)

func (l *Logger) LogMode(lev logger.LogLevel) logger.Interface {
	return &Logger{}
}
func (l *Logger) Info(ctx context.Context, msg string, data ...interface{}) {
	logx.WithContext(ctx).Infof(msg, data)
}
func (l *Logger) Warn(ctx context.Context, msg string, data ...interface{}) {
	logx.WithContext(ctx).Errorf(msg, data)
}
func (l *Logger) Error(ctx context.Context, msg string, data ...interface{}) {
	logx.WithContext(ctx).Errorf(msg, data)
}
func (l *Logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)

	// 1. 发生错误
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// RecordNotFound 通常作为业务逻辑判断，用 Infow 记录
			logx.WithContext(ctx).Infow("Database Record Not Found", l.buildFields(elapsed, fc)...)
		} else {
			// 真正的数据库异常，必须用 Errorw
			fields := l.buildFields(elapsed, fc)
			fields = append(fields, logx.Field("catch_error", err.Error()))
			logx.WithContext(ctx).Errorw("Database Error", fields...)
		}
		return
	}

	// 2. 触发慢 SQL 阈值
	if elapsed > slowThreshold {
		logx.WithContext(ctx).Sloww("Database Slow Query", l.buildFields(elapsed, fc)...)
		return
	}

	// 3. 正常查询：不要在这里用 if 拦死！直接交给 go-zero 的底层去判断当前环境是否允许打印 Info
	// 如果线上关闭了 Info 级别，go-zero 内部会直接 return，Fc() 也不会被真正执行，性能损耗极小
	logx.WithContext(ctx).Infow("Database Query", l.buildFields(elapsed, fc)...)
}

// 辅助方法：将字段组装抽离
func (l *Logger) buildFields(elapsed time.Duration, fc func() (string, int64)) []logx.LogField {
	sql, rows := fc()

	// 耗时转换为 float64 毫秒，方便监控大盘直接计算 AVG / P99
	costMs := float64(elapsed.Nanoseconds()) / 1e6

	fields := []logx.LogField{
		logx.Field("db.statement", sql),
		logx.Field("db.duration_ms", costMs),
		logx.Field("db.rows_affected", rows),
	}

	// 自动获取是哪个业务函数、哪一行代码触发的这条 SQL
	// utils.FileWithLineNum 是 GORM 内部的高效实现，会自动跳过 GORM 自身的源码层
	if caller := utils.FileWithLineNum(); caller != "" {
		fields = append(fields, logx.Field("db.caller", caller))
	}

	return fields
}
