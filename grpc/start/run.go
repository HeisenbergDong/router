package start

import (
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"router/global"
	"router/grpc/controller"
	"router/grpc/pb"
	"router/utils"
)

type GrpcServer struct {
	grpcAddress string
}

func NewGrpcServer() *GrpcServer {
	return &GrpcServer{
		grpcAddress: global.CONFIG.Grpc.JwtAddress,
	}
}

func (g *GrpcServer) Run() {
	go func() {
		defer utils.RecoverPanic()
		listen, err := net.Listen("tcp", g.grpcAddress)
		if err != nil {
			global.LOG.Error("grpc server", zap.Any(" Failed to listen:", err))
		}
		s := grpc.NewServer()
		// 服务注册
		pb.RegisterJWTServer(s, &controller.JwtController{})

		global.LOG.Info("Listen on " + g.grpcAddress)

		if err := s.Serve(listen); err != nil {
			global.LOG.Error("grpc server", zap.Any("Failed to serve:", err))
		}
	}()
}
