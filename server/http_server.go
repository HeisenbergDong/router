package server

import (
	"go.uber.org/zap"
	"net/http"
	"router/global"
	"strconv"
	"time"
)

type GatewayServer struct {
	timeout     time.Duration
	host        string
	port        int
	contextPath string
}

func NewGatewayServer() *GatewayServer {
	return &GatewayServer{
		timeout:     time.Second * global.CONFIG.Server.Timeout,
		host:        global.CONFIG.Server.Host,
		port:        global.CONFIG.Server.Port,
		contextPath: global.CONFIG.Server.ContextPath,
	}
}

func (r *GatewayServer) Start() {
	gatewayProxy := NewGatewayProxy()
	http.HandleFunc(r.contextPath, gatewayProxy.dispatch)
	http.Handle("/favicon.ico", http.FileServer(http.Dir("static")))

	// 服务监听地址
	addr := r.host + ":" + strconv.Itoa(r.port)
	global.LOG.Info("服务启动中，服务绑定端口号：", zap.Any("port", r.port))
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		global.LOG.Fatal("API 网关服务启动失败", zap.Any("err", err))
	}
}