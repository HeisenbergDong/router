package filter

import (
	"go.uber.org/zap"
	"router/context"
	"router/global"
	"sync"
)

var afterResponseFunc []Handler
var afterResponseLock = &sync.RWMutex{}

func registerAfterResponse(handle Handler) {
	afterResponseLock.Lock()
	defer afterResponseLock.Unlock()
	global.LOG.Info("注册响应后过滤器，过滤器名称是：", zap.Any("handleName",handle.Name))

	result := make([]Handler, len(afterResponseFunc) + 1)

	if len(afterResponseFunc) == 0 {
		result = append(afterResponseFunc, handle)
	} else {
		for idx, h := range afterResponseFunc {
			if h.Priority > handle.Priority {
				if idx == 0 {
					// 第一个元素
					f := []Handler{handle}
					result = append(f, afterResponseFunc[0])
				} else if idx + 1 == len(afterResponseFunc) {
					// 最后一个元素
					last := afterResponseFunc[idx]
					v := append(afterResponseFunc[:idx], handle)
					result = append(v, last)
				} else {
					// 中间元素
					v := append(afterResponseFunc[:idx], handle)
					result = append(v, afterResponseFunc[idx:]...)
				}
				break
			}
			result = append(afterResponseFunc, handle)
		}
	}
	afterResponseFunc = result
}

func AfterResponseFilter(ctx *context.GatewayContext)  error {
	afterResponseLock.RLock()
	defer afterResponseLock.RUnlock()

	for _, f := range afterResponseFunc {
		 err := f.Handle(ctx)
		if err != nil {
			global.LOG.Error("AfterResponseFilter: ",zap.Any("err",err))
			return  err
		}
	}

	return  nil
}
