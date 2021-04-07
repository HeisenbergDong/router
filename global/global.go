package global

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"router/config"
)

var (
	VIPER   *viper.Viper
	CONFIG  config.Config
	ADDRESS string
	ROUTER  *gin.Engine
	LOG     *zap.Logger
	DB      *gorm.DB
	REDIS   *redis.Client
)
