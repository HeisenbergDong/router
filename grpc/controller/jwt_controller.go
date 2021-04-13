package controller

import (
	"context"
	"router/grpc/pb"
	"router/grpc/service"
)

type JwtController struct{}

func (jc *JwtController) CreateToken(ctx context.Context, user *pb.SysUser) (*pb.TokenMessage, error) {
	return service.CreateToken(user)
}
