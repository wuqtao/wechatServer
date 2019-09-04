package redisClient

import (
	"github.com/dbldqt/wechatServer/config"
	"github.com/gomodule/redigo/redis"
	"log"
	"time"
)

var redisClient *redis.Pool

func GetRedisClient() *redis.Pool{
	if redisClient != nil{
		return redisClient
	}else{
		redisConf := config.GetConfig().Redis
		redisClient = &redis.Pool{
				IdleTimeout: 10 * time.Second,
				Wait:        true,
				Dial: func() (redis.Conn, error) {
					con, err := redis.Dial("tcp", redisConf.Host,
						redis.DialPassword(redisConf.Password),
						redis.DialDatabase(redisConf.Database),
						redis.DialConnectTimeout(5*time.Second),
						redis.DialReadTimeout(5*time.Second),
						redis.DialWriteTimeout(5*time.Second))
					if err != nil {
						log.Println("connect ")
						return nil, err
					}
					return con, nil
				},
			}
	}
	return redisClient
}

func GetConn() redis.Conn{
	return GetRedisClient().Get()
}

func Lock(keyStr string,lockNum int,expire int) bool{
	conn := redisClient.Get()
	defer conn.Close()
	result,err := conn.Do("set",keyStr,lockNum,"nx","ex",expire)
	if err != nil{
		log.Println("redisClient set error "+err.Error())
		return false
	}
	re,ok := result.(string);if ok{
		if re == "OK"{
			return true
		}
		return false
	}
	return false
}

func UnLock(keyStr string,lockNum int) bool{
	conn := redisClient.Get()
	defer conn.Close()
	num,err := redis.Int(conn.Do("get",keyStr))
	if err != nil{
		return false
	}
	if num == lockNum{
		_,err = conn.Do("del",keyStr)
		if err != nil{
			log.Println("del lock key error "+err.Error())
			return false
		}
		return true
	}
	log.Println("try to unlock other lock")
	return false
}

func IsLocked(keyStr string) bool{
	conn := redisClient.Get()
	defer conn.Close()
	result,err := redis.Int(conn.Do("get",keyStr))
	if err != nil{
		log.Println("check redis lock error")
		return true
	}
	if result != 0{
		return true
	}
	return false
}