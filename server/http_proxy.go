package server

import (
	"go.uber.org/zap"
	"net/http"
	"net/http/httputil"
	"router/context"
	"router/filter"
	"router/global"
	"router/router"
)

type GatewayProxy struct{}

func NewGatewayProxy() *GatewayProxy {
	return &GatewayProxy{}
}

// dispatch API 请求分发
func (r *GatewayProxy) dispatch(w http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(w, req)

	defer func() {
		if err := recover(); err != nil {
			global.LOG.Error("dispatch:", zap.Any("err", err))
			r.globalRecover(ctx, err)
		}
	}()

	// 收到消息之后，统一进行处理
	need := router.IsFilter(ctx.Request.URL.Path)
	err := filter.BeforeRequestFilter(need, ctx)
	if err != nil {
		ErrorHandle(ctx, err)
		return
	}

	// 请求后端服务
	err = r.httpProxy(ctx)
	if err != nil {
		ErrorHandle(ctx, err)
		return
	}

	// 执行后置过滤器
	err = filter.AfterResponseFilter(ctx)
	if err != nil {
		ErrorHandle(ctx, err)
		return
	}
}

// httpProxy 发起 http 请求
func (r *GatewayProxy) httpProxy(ctx *context.GatewayContext) error {

	// 匹配路由
	path := ctx.Request.URL.Path
	remoteUrl, route, err := router.Match(path)
	if err != nil {
		global.LOG.Error("match route is fail: ", zap.Any("err", err))
		return NewError(ProxyUrlNotFound, err.Error())
	}
	ctx.RemoteURL = remoteUrl

	// 创建代理对象
	proxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.Header.Add("X-Forwarded-Host", req.Host)
			req.Header.Add("X-Origin-Host", remoteUrl.Host)
			req.Host = remoteUrl.Host
			req.URL.Scheme = remoteUrl.Scheme
			req.URL.Host = remoteUrl.Host
			req.URL.Path = remoteUrl.Path

			r.filterSensitiveHeaders(req)

			if ctx.Request.URL.RawQuery == "" || req.URL.RawQuery == "" {
				req.URL.RawQuery = ctx.Request.URL.RawQuery + req.URL.RawQuery
			} else {
				req.URL.RawQuery = ctx.Request.URL.RawQuery + "&" + req.URL.RawQuery
			}
			ctx.Request = req
			global.LOG.Debug("director: ", zap.Any("remoteUrl", remoteUrl))
		},
		ModifyResponse: func(resp *http.Response) error {
			global.LOG.Debug("modify response:", zap.Any("remoteUrl", remoteUrl))
			ctx.Response = resp
			err := filter.BeforeResponseFilter(ctx)
			if err != nil {
				global.LOG.Info("BeforeResponse Stop")
			}
			return nil
		},
		ErrorHandler: r.ErrorHandler,
		Transport:    defaultGatewayTransport.GetTransport(route),
	}

	proxy.ServeHTTP(ctx.ResponseWriter, ctx.Request)

	global.LOG.Info("Http Proxy", zap.Any("请求完成, 请求地址：", path), zap.Any("目标地址：", remoteUrl))
	return nil
}

// filterSensitiveHeaders 过滤掉请求 Header 中配置的 Key
func (r *GatewayProxy) filterSensitiveHeaders(req *http.Request) {
	for _, header := range global.CONFIG.GatewayRouter.SensitiveHeaders {
		req.Header.Del(header)
	}
}

// globalErrorHandle Proxy 代理处理错误
func (r *GatewayProxy) ErrorHandler(w http.ResponseWriter, request *http.Request, err error) {
	global.LOG.Error(err.Error())
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusBadGateway)
	w.Write([]byte(NewError(ProxyError, err.Error()).Error()))
}

// globalRecover 捕获全局异常
func (r *GatewayProxy) globalRecover(ctx *context.GatewayContext, errMsg interface{}) {
	if ctx.ResponseWriter != nil {
		ctx.ResponseWriter.Header().Set("Content-Type", "application/json;charset=UTF-8")
		ctx.ResponseWriter.WriteHeader(http.StatusBadGateway)
		ctx.ResponseWriter.Write([]byte(NewError(ProxyError, errMsg).Error()))
	}
}
