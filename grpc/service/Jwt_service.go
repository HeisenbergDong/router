package service

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"
	"router/global"
	"router/grpc/dao"
	"router/grpc/model"
	"router/grpc/pb"
	"time"
)

func IsBlacklist(jwt string) bool {
	isNotFound := dao.IsBlacklist(jwt)
	return !isNotFound
}

func JsonInBlacklist(jwtList model.JwtBlacklist) (err error) {
	if err = dao.JsonInBlacklist(jwtList); err != nil {
		global.LOG.Error("JsonInBlacklist is fail", zap.Any("err", err))
	}
	return err
}

func GetRedisJWT(userName string) (err error, redisJWT string) {
	redisJWT, err = global.REDIS.Get(userName).Result()
	return err, redisJWT
}

func SetRedisJWT(jwt string, userName string) (err error) {
	// 此处过期时间等于jwt过期时间
	timer := time.Duration(global.CONFIG.JWT.ExpiresTime) * time.Second
	err = global.REDIS.Set(userName, jwt, timer).Err()
	return err
}

//func GetUserID(c *gin.Context) uint {
//	if claims, exists := c.Get("claims"); !exists {
//		global.LOG.Error("从Gin的Context中获取从jwt解析出来的用户ID失败, 请检查路由是否使用jwt中间件")
//		return 0
//	} else {
//		waitUse := claims.(*model.CustomClaims)
//		return waitUse.ID
//	}
//}

//登录以后签发jwt
func CreateToken(user *pb.SysUser) (*pb.TokenMessage, error) {
	j := &JWT{SigningKey: []byte(global.CONFIG.JWT.SigningKey)} // 唯一签名
	claims := model.CustomClaims{
		UUID:        user.Uuid,
		ID:          user.Id,
		NickName:    user.NickName,
		UserName:    user.UserName,
		AuthorityId: user.AuthorityId,
		BufferTime:  global.CONFIG.JWT.BufferTime, // 缓冲时间1天 缓冲时间内会获得新的token刷新令牌 此时一个用户会存在两个有效令牌 但是前端只留一个 另一个会丢失
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 1000,                          // 签名生效时间
			ExpiresAt: time.Now().Unix() + global.CONFIG.JWT.ExpiresTime, // 过期时间 7天  配置文件
			Issuer:    "qmPlus",                                          // 签名的发行者
		},
	}
	token, err := j.CreateToken(claims)
	if err != nil {
		global.LOG.Error("获取token失败", zap.Any("err", err))
		return &pb.TokenMessage{}, errors.New("create token fail")
	}
	return &pb.TokenMessage{Token: token, ExpiresAt: claims.StandardClaims.ExpiresAt * 1000}, nil
}

type JWT struct {
	SigningKey []byte
}

func NewJWT() *JWT {
	return &JWT{SigningKey: []byte(global.CONFIG.JWT.SigningKey)}
}

// 创建一个token
func (j *JWT) CreateToken(claims model.CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// 解析 token
func (j *JWT) ParseToken(tokenString string) (*model.CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &model.CustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*model.CustomClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, TokenInvalid
	} else {
		return nil, TokenInvalid
	}

}

var (
	TokenExpired     = errors.New("token is expired")
	TokenNotValidYet = errors.New("token not active yet")
	TokenMalformed   = errors.New("that's not even a token")
	TokenInvalid     = errors.New("couldn't handle this token")
)
