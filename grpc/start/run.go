package start

import (
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"router/global"
	"router/grpc/controller"
	"router/grpc/pb"
)

type GrpcServer struct {
	grpcAddress string
}

func NewGrpcServer() *GrpcServer {
	return &GrpcServer{
		grpcAddress: global.CONFIG.Server.GrpcAddress,
	}
}

func (g *GrpcServer) Run() {
	go func() {
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

//
//const (
//	Address = "127.0.0.1:9528"
//)
//
//func main() {
//	listen, err := net.Listen("tcp", Address)
//	if err != nil {
//		log.Fatalf("Failed to listen: %v", err)
//	}
//	s := grpc.NewServer()
//	// 服务注册
//	pb.RegisterJWTServer(s, &controller.JwtController{})
//
//	log.Println("Listen on " + Address)
//
//	if err := s.Serve(listen); err != nil {
//		log.Fatalf("Failed to serve: %v", err)
//	}
//}
