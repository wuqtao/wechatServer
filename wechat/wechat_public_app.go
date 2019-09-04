package wechat

import (
	"errors"
	"github.com/dbldqt/wechatSDK/mp/accesstoken"
	"github.com/dbldqt/wechatServer/config"
	"github.com/dbldqt/wechatServer/redisClient"
	"github.com/gomodule/redigo/redis"
	"log"
	"math/rand"
	"strconv"
	"time"
)

const (
	UpdateTokenLockKey = "updateTokenLock"
	TokenHashKey = "token"
	TokenExpireHashKey = "expireAt"
)
var wechatPulicApp *WechatPublicApp

type WechatPublicApp struct {
	config *config.WechatPublicConf
}

func NewWechatPublicApp(config *config.WechatPublicConf) *WechatPublicApp{
	wechatPulicApp = &WechatPublicApp{
		config:config,
	}
	return wechatPulicApp
}

func GetWechatPublicApp() *WechatPublicApp{
	return wechatPulicApp
}

//使用一个redis hash结构存储微信的accessToken,key accesstoken val expireAt val isLock bool lockedAt val
func (self *WechatPublicApp)getAccessTokenKey() string{
	return "accessToken"+self.config.AppID
}

func (self *WechatPublicApp)GetAccessToken() (string,int64,error){
	//先查询是否被锁定
	conn := redisClient.GetConn()
	defer conn.Close()
	retMap,err := redis.StringMap(conn.Do("hgetall",self.getAccessTokenKey()))
	if err != nil{
		log.Println("query accessToken from redis error "+err.Error())
		return "",0,err
	}
	token,ok := retMap[TokenHashKey];if !ok{
		log.Println("query accessToken from redis token not exist")
		return "",0,errors.New("query accessToken from redis token not exist")
	}
	expireAt,ok := retMap[TokenExpireHashKey];if !ok{
		log.Println("query accessToken from redis expireAt not exist")
		return "",0,errors.New("query accessToken from redis expireAt not exist")
	}
	expireAtInt,err := strconv.ParseInt(expireAt,10,64)
	if err != nil{
		log.Println("parse expire error")
		return token,0,err
	}
	return token,expireAtInt,nil
}

func (self *WechatPublicApp)Start(){
	go func(){
		for{
			self.checkAccessTokenExpire()
			time.Sleep(time.Duration(self.config.LoopTime)*time.Second)
		}
	}()
}

func (self *WechatPublicApp)checkAccessTokenExpire(){
	token,expireAt,err := self.GetAccessToken()
	if err != nil{
		log.Println("check accessToken error "+err.Error())
		self.UpdateAccessToken()
		return
	}
	if token == "" || time.Now().Unix() >= expireAt{
		self.UpdateAccessToken()
	}
}

func (self *WechatPublicApp)setAccessToken(){
	token,expireIn,err := accesstoken.GetAccessToken(self.config.AppID,self.config.AppSecret)
	if err != nil{
		log.Println("request accesstoken error "+err.Error())
		return
	}
	conn := redisClient.GetConn()
	defer conn.Close()
	expireAt := time.Now().Unix()+int64(expireIn)-self.config.AheadTimeTime
	_,err = conn.Do("hmset",self.getAccessTokenKey(),TokenHashKey,token,TokenExpireHashKey,expireAt)
	if err != nil{
		log.Println("set accessToken error "+err.Error())
		return
	}
}

func (self *WechatPublicApp)UpdateAccessToken(){
	rand.Seed(time.Now().Unix())
	lockNum := rand.Int()
	//上锁再进行更新,锁过期时间2分钟
	if !redisClient.Lock(UpdateTokenLockKey,lockNum,120){
		log.Println("lock update token lock failed")
		return
	}
	self.setAccessToken()
	redisClient.UnLock(UpdateTokenLockKey,lockNum)
}


