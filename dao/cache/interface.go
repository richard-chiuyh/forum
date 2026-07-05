package cache

import (
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type Manager struct {
	Post PostCache
}

func NewManager(conf redis.RedisConf) (Manager, error) {
	rdsCli, err := redis.NewRedis(conf)
	if err != nil || !rdsCli.Ping() {
		logx.Severef("redis connection err(%+v)", err)
		return Manager{}, err
	}
	return Manager{
		Post: NewPostOpr(rdsCli, 60*60*24), // 过期时间：24小时
	}, nil
}

type PostCache interface {
}
