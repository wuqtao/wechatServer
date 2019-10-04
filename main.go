package main

import (
	"flag"
	"fmt"
	"github.com/dbldqt/wechatServer/config"
	"github.com/dbldqt/wechatServer/controller"
	"github.com/dbldqt/wechatServer/wechat"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
	"github.com/kataras/iris/mvc"
	"log"
	"os"
	"strconv"
)
var configFile string
var conf *config.Config
var test bool

func init(){
	flag.StringVar(&configFile,"conf","./config.toml","assign the config file path")
	flag.BoolVar(&test,"test",false,"is test config file")
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func main(){
	flag.Parse()

	//测试配置文件是否正确
	if test{
		err := config.TestConfig(configFile)
		if err != nil{
			fmt.Println(err.Error())
			return
		}
		fmt.Println("config file is right")
		return
	}

	conf,err :=config.ParseConfig(configFile)
	if err != nil{
		fmt.Println("parse config error "+err.Error())
		return
	}

	//初始化log
	file,err := os.OpenFile(conf.Server.LogFile,os.O_CREATE | os.O_APPEND | os.O_RDWR,0666)
	if err != nil{
		log.Println("open log file error "+err.Error())
	}else{
		log.SetOutput(file)
	}

	wechat.NewWechatPublicApp(conf.WechatPublic).Start()

	app := iris.Default()
	app.Use(recover.New())
	app.Use(logger.New())

	mvc.Configure(app.Party("accesstoke"), func(app *mvc.Application) {
		app.Handle(new(controller.AccessToken))
		app.Handle(new(controller.Login))
	})
	app.Run(iris.Addr(":"+strconv.Itoa(conf.Server.Port)))
}
