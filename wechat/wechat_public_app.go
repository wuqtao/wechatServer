package wechat

import (
	"github.com/dbldqt/wechatServer/config"
	"github.com/dbldqt/wechatSDK/mp/accesstoken"
	"github.com/dbldqt/wechatServer/redis"
	"log"
	"strconv"
	"time"
)
var wechatPulicApp *WechatPublicApp

type WechatPublicApp struct {
	config *config.WechatPublicConf
}

func NewWechatPublicApp(config *config.WechatPublicConf) *WechatPublicApp{

}

//使用一个redis hash结构存储微信的accessToken,key accesstoken val expireAt val isLock bool lockedAt val
func (self *WechatPublicApp)getAccessTokenKey() string{
	return "accessTokenkey"+self.config.AppID
}

func (self *WechatPublicApp)GetAccessToken() string{
	result,err := redis.GetRedisClient().HGetAll(self.getAccessTokenKey()).Result()
	if err != nil{
		log.Println("query accesstoken from redis error "+err.Error())
		return ""
	}
	token,ok := result["accessToken"]
	if ok{
		return token
	}
	return ""
}

func (self *WechatPublicApp)Start(){
	go func(){
		for{
			result,err := redis.GetRedisClient().HGetAll(self.getAccessTokenKey()).Result()
			if err != nil{
				log.Println("query accesstoken from redis error "+err.Error())
				continue
			}
			duration,_ := time.ParseDuration(strconv.Itoa(self.config.LoopTime))
			expiredAt,ok := result["expiredAt"]
			if !ok {
				self.UpdateAccessToken()
			}
			expireInt64,_ := strconv.ParseInt(expiredAt,10,64)
			if time.Unix(expireInt64,0).Before(time.Now()){
				self.UpdateAccessToken()
			}
			time.Sleep(duration)
		}

	}()
}
func (self *WechatPublicApp)checkAccessTokenExpire(){
	//先查询是否被锁定
	result,err := redis.GetRedisClient().HGetAll(self.getAccessTokenKey()).Result()
	if err != nil{
		log.Println("query accesstoken from redis error "+err.Error())
		return
	}
	token,ok := result["token"]
	isLook,ok1 := result["isLock"]

	if ok{

	}
}
func (self *WechatPublicApp)UpdateAccessToken(){
	//上锁再进行更新
	_,err := redis.GetRedisClient().HMSet(self.getAccessTokenKey(),map[string]interface{}{
		"isLock":true,
		"lockedAt":time.Now().Unix(),
	}).Result()
	if err != nil{
		log.Println("lock redis error "+err.Error())
		return
	}
	token,expireIn,err := accesstoken.GetAccessToken(self.config.AppID,self.config.AppSecret)
	if err != nil{
		log.Println("request accesstoken error "+err.Error())
		return
	}

	expireAt := time.Now().Unix()+int64(expireIn)-self.config.AheadTimeTime
	redis.GetRedisClient().HMSet(self.getAccessTokenKey(),map[string]interface{}{
		"token":token,								//存储accessToken
		"expiredAt":expireAt,						//token过期时间
		"isLock":false,								//用于分布式并发更新控制
		"lockedAt":0,								//锁定时间，1分钟不更新则任务超时，其他服务可以重新锁定并更新accessToken
	})
}


