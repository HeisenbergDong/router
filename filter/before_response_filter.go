package filter

import (
	"go.uber.org/zap"
	"router/context"
	"router/global"
	"sync"
)

var beforeResponseFunc []Handler
var beforeResponseLock = &sync.RWMutex{}

func registerBeforeResponse(handle Handler) {
	beforeResponseLock.Lock()
	defer beforeResponseLock.Unlock()
	global.LOG.Info("注册响应前过滤器，过滤器名称是：", zap.Any("handleName: ",handle.Name))

	result := make([]Handler, len(afterResponseFunc) + 1)

	if len(beforeResponseFunc) == 0 {
		result = append(beforeResponseFunc, handle)
	} else {
		for idx, h := range beforeResponseFunc {
			if h.Priority > handle.Priority {
				if idx == 0 {
					// 第一个元素
					f := []Handler{handle}
					result = append(f, beforeResponseFunc[0])
				} else if idx + 1 == len(beforeResponseFunc) {
					// 最后一个元素
					last := beforeResponseFunc[idx]
					v := append(beforeResponseFunc[:idx], handle)
					result = append(v, last)
				} else {
					// 中间元素
					v := append(beforeResponseFunc[:idx], handle)
					result = append(v, beforeResponseFunc[idx:]...)
				}
				break
			}
			result = append(beforeResponseFunc, handle)
		}
	}
	beforeResponseFunc = result
}

func BeforeResponseFilter(ctx *context.GatewayContext) error {
	beforeResponseLock.RLock()
	defer beforeResponseLock.RUnlock()

	for _, f := range beforeResponseFunc {
		err := f.Handle(ctx)
		if err != nil {
			global.LOG.Error("BeforeResponseFilter: ",zap.Any("err",err))
			return  err
		}
	}
	return nil
}