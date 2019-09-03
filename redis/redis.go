package redis

import (
	"github.com/dbldqt/wechatServer/config"
	"github.com/go-redis/redis"
)

var redisClient *redis.Client

func GetRedisClient() *redis.Client{
	if redisClient != nil{
		return redisClient
	}else{
		redisConf := config.GetConfig().Redis
		redisClient = redis.NewClient(&redis.Options{
			Addr:redisConf.Host,
			Password:redisConf.Password,
			DB:redisConf.Database,
		})
		return redisClient
	}
}
