package config

import (
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"log"
	"sync"
)
var config *Config

func GetConfig() *Config{
	if config != nil{
		return config
	}
	log.Fatalln("config not init")
	return nil
}

type ServerConf struct {
	LogFile string		`toml:"logFile"`
	Port	int			`toml:"port"`
}

type RedisConf struct {
	Host string 		`toml:"host"`
	Password string		`toml:"password"`
	Database int		`toml:"database"`
}

type SecurityConf struct {
	ApiSignKey string	`toml:"apiSignKey"`
	AdminIpList []string 	`toml:"adminIpList"`
	AdminToken	string		`toml:"adminToken"`
	UseIpWhiteList bool		`toml:"useIpWhiteList"`
	IpList []string			`toml:"ipList"`
}

type WechatPublicConf struct {
	AppID string		`toml:"appID"`
	AppSecret string	`toml:"appSecret"`
	Token string		`toml:"token"`
	NotifyUrl []string	`toml:"notifyUrl"`
	EncodeAesKey string	`toml:"encodeAesKey"`

	AheadTimeTime int64	`toml:"aheadTime"`
	LoopTime int	`toml:"loopTime"`
}

type MiniProgramConf struct {

}

type WechatWebConf struct {

}

type Config struct {
	sync.RWMutex
	Server *ServerConf 				`toml:"server"`
	Redis *RedisConf				`toml:"redis"`
	Security *SecurityConf		`toml:"security"`
	WechatPublic *WechatPublicConf	`toml:"wechatPublic"`
	MiniProgram *MiniProgramConf	`toml:"miniProgram"`
	WechatWeb *WechatWebConf		`toml:"wechatWeb"`
}
/**
	解析配置文件
 */
func ParseConfig(filePath string) (*Config,error){
	fileData,err := ioutil.ReadFile(filePath)
	if err != nil{
		return nil,err
	}

	if _,err = toml.Decode(string(fileData),&config);err != nil{
		return nil,err
	}
	return config,nil
}
/**
	用于测试配置文件
 */
func TestConfig(filePath string) error{
	var tmpConfig Config
	if _,err := toml.DecodeFile(string(filePath),&tmpConfig);err != nil{
		log.Fatalln("decode config error "+err.Error())
		return err
	}
	return nil
}

