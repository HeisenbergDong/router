package authcontroller

import (
	"context"
	"google.golang.org/grpc"
	"router/global"
	"router/grpc/pb"
)

func CasController(sub string, obj string, act string) (success bool, err error) {
	conn, err := grpc.Dial(global.CONFIG.Grpc.AuthAddress, grpc.WithInsecure())
	if err != nil {
		global.LOG.Info(err.Error())
	}
	defer conn.Close()
	// 初始化客户端
	c := pb.NewAUTHClient(conn)
	req := &pb.Req{
		Sub: sub,
		Obj: obj,
		Act: act,
	}
	rep, err := c.IsAuth(context.Background(), req)
	return rep.Success, err
}
