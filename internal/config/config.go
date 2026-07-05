package config

import (
	"forum/dbcore"
	"forum/escore"

	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	// SentryConfig monitoring.SentryConfig
	MysqlConfig dbcore.MysqlConfig
	ES          escore.ESConfig
	Kafka       struct {
		Address []string
	}
}
