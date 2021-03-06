package main

import (
	"router/global"
	"router/grpc/start"
	"router/middleware"
	"router/pkg/gorm"
	"router/pkg/redis"
	"router/pkg/viper"
	"router/pkg/zap"
	"router/router"
	"router/server"
)

func main() {
	global.VIPER = viper.InitViper() // 初始化Viper配置文件
	global.LOG = zap.InitZap()       // 初始化zap日志库
	global.DB = gorm.InitGorm()      // 初始化gorm连接数据库
	global.REDIS = redis.InitRedis() // 初始化redis服务
	router.InitGatewayRouter()       // 初始化gateway
	middleware.JWTRegisterFilter()   // 注册JWT拦截
	gatewayServer := server.NewGatewayServer()
	gatewayServer.Start() // 启动Gateway服务
	grpcServer := start.NewGrpcServer()
	grpcServer.Run() // 启动GRPC服务
	select {}
}
