package cache

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type postOpr struct {
	rdsCli      *redis.Redis
	cacheExpire int // 缓存过期时间
}

func NewPostOpr(cli *redis.Redis, expire int) PostCache {
	return &postOpr{
		rdsCli:      cli,
		cacheExpire: expire,
	}
}
