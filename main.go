package main

import (
	"router/global"
	"router/pkg/gorm"
	"router/pkg/redis"
	"router/pkg/viper"
	"router/pkg/zap"
	"router/router"
	"router/server"
)

func main()  {
	global.VIPER = viper.InitViper()     		 // 初始化Viper配置文件
	global.LOG = zap.InitZap()           		 // 初始化zap日志库
	global.DB = gorm.InitGorm()          		 // 初始化gorm连接数据库
	global.REDIS = redis.InitRedis()     		 // 初始化redis服务
	router.InitGatewayRouter()					 // 初始化gateway
	gatewayServer := server.NewGatewayServer()
	gatewayServer.Start()						 // 启动服务
}