package main

import (
	"google.golang.org/grpc"
	"log"
	"net"
	"router/grpc/controller"
	"router/grpc/pb"
)

const (
	Address = "127.0.0.1:9528"
)

func main() {
	listen, err := net.Listen("tcp", Address)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	// 服务注册
	pb.RegisterJWTServer(s, &controller.JwtController{})

	log.Println("Listen on " + Address)

	if err := s.Serve(listen); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
