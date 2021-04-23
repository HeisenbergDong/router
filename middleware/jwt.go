package middleware

import (
	"errors"
	"go.uber.org/zap"
	"router/context"
	"router/filter"
	"router/global"
	"router/grpc/client/authcontroller"
	"router/grpc/model"
	"router/grpc/service"
	"strconv"
	"time"
)

func JWTRegisterFilter() {
	filter.RegisterFilter(filter.BeforeRequest, filter.Handler{
		Name:     "JWT",
		Priority: 1,
		Handle: func(ctx *context.GatewayContext) error {
			// 我们这里jwt鉴权取头部信息 x-token 登录时回返回token信息 这里前端需要把token存储到cookie或者本地localStorage中 不过需要跟后端协商过期时间 可以约定刷新令牌或者重新登录
			token := ctx.Request.Header.Get("x-token")
			if token == "" {
				return errors.New("token is null")
			}
			if service.IsBlacklist(token) {
				return errors.New("token in black list")
			}
			j := service.NewJWT()
			// parseToken 解析token包含的信息
			claims, err := j.ParseToken(token)
			if err != nil {
				if err == service.TokenExpired {
					return errors.New("token is expired")
				}
				return errors.New("parse token fail")
			}
			if claims.ExpiresAt-time.Now().Unix() < claims.BufferTime {
				claims.ExpiresAt = time.Now().Unix() + global.CONFIG.JWT.ExpiresTime
				newToken, _ := j.CreateToken(*claims)
				newClaims, _ := j.ParseToken(newToken)
				ctx.Request.Header.Set("new-token", newToken)
				ctx.Request.Header.Set("new-expires-at", strconv.FormatInt(newClaims.ExpiresAt, 10))
				if global.CONFIG.System.UseMultipoint {
					err, RedisJwtToken := service.GetRedisJWT(newClaims.UserName)
					if err != nil {
						global.LOG.Error("get redis jwt failed", zap.Any("err", err))
					} else { // 当之前的取成功时才进行拉黑操作
						_ = service.JsonInBlacklist(model.JwtBlacklist{Jwt: RedisJwtToken})
					}
					// 无论如何都要记录当前的活跃状态
					_ = service.SetRedisJWT(newToken, newClaims.UserName)
				}
			}
			//CASbin
			sub := claims.AuthorityId
			act := ctx.Request.Method
			obj := ctx.Request.URL.RequestURI()
			success, err := authcontroller.CasController(sub, act, obj)
			if global.CONFIG.System.Env == "develop" || success {
				return nil
			} else {
				return err
			}
		},
	})
}
