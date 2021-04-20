package controller

import (
	"context"
	"router/grpc/model"
	"router/grpc/pb"
	"router/grpc/service"
)

type JwtController struct{}

func (jc *JwtController) CreateToken(ctx context.Context, user *pb.SysUser) (tokenMessage *pb.TokenMessage, err error) {
	tokenMessage, err = service.CreateToken(user)
	return tokenMessage, err
}

func (jc *JwtController) GetRedisJWT(ctx context.Context, req *pb.GetRedisJWTReq) (getRedisRep *pb.GetRedisRep, err error) {
	err, redisJWT := service.GetRedisJWT(req.UserName)
	getRedisRep.RedisJwt = redisJWT
	return getRedisRep, err
}

func (jc *JwtController) SetRedisJWT(ctx context.Context, req *pb.SetRedisJWTReq) (rep *pb.SetRedisJWTRep, err error) {
	return &pb.SetRedisJWTRep{}, service.SetRedisJWT(req.Token, req.UserName)
}

func (jc *JwtController) JsonInBlacklist(ctx context.Context, req *pb.JsonInBlacklistReq) (rep *pb.JsonInBlacklistRep, err error) {
	var jwtList model.JwtBlacklist
	jwtList.Jwt = req.BlackJWT
	return &pb.JsonInBlacklistRep{}, service.JsonInBlacklist(jwtList)
}

func (jc *JwtController) GetUserID(ctx context.Context, req *pb.GetUserIDReq) (rep *pb.GetUserIDRep, err error) {
	id, err := service.GetUserID(req.Token)
	return &pb.GetUserIDRep{Id: id}, err
}
