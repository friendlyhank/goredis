package rds

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"strconv"
	"strings"
)

var defaultServer string

// Init -
func Init() {
	logs.Debug("|foundation|init|rds|Init")

	redissource := beego.AppConfig.DefaultString("redis","127.0.0.1:6379,150,123456")
	if rdss := strings.Split(redissource, ","); len(rdss) == 3 {
		// address,connect,password
		maxIdle, _ := strconv.Atoi(rdss[1])
		InitRedisServer(rdss[0], rdss[2], maxIdle)
	}
}

// InitRedisServer -
func InitRedisServer(server, password string, maxIdle int) {
	defaultServer = server
	AddRedisServer(server, password, maxIdle)
}

// GetRedisDefault -
func GetRedisDefault() *RedisSource {
	if len(defaultServer) == 0 {
		logs.Debug("|foundation|rds|redis|GetRedisDefault|len(defaultServer)==0|redisSourceMap|%+v", redisSourceMap)
		for _, s := range redisSourceMap {
			return s
		}
	}
	return redisSourceMap[defaultServer]
}
