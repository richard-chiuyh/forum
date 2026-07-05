package svc

import (
	"forum/dao/cache"
	"forum/dao/db"
	"forum/dao/es"
	"forum/internal/config"

	"github.com/zeromicro/go-zero/core/logx"
)

type ServiceContext struct {
	Config    config.Config
	DBManager db.Manager
	// Kafka        broker.Broker
	CacheManager cache.Manager
	ESManager    es.Manager
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化数据库
	dbManager, err := db.NewManager(&c.MysqlConfig)
	if err != nil {
		logx.Severef("db connection err(%+v)", err)
		panic(0)
	}
	// 初始化缓存
	redisManager, err := cache.NewManager(c.Redis.RedisConf)
	if err != nil {
		logx.Severef("redis connection err(%+v)", err)
		panic(0)
	}
	// 初始化 ES 客户端
	esManager, err := es.NewManager(c)
	if err != nil {
		logx.Severef("es connection err(%+v)", err)
		panic(0)
	}

	return &ServiceContext{
		Config:    c,
		DBManager: dbManager,
		// Kafka:        bk,
		CacheManager: redisManager,
		ESManager:    esManager,
	}
}
