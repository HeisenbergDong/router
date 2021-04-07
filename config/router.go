package config

import "time"

type Routers struct {
	Routers          map[string]RouterDetails
	IgnoredPatterns  []string `mapstructure:"ignoredPatterns" json:"ignoredPatterns" yaml:"ignoredPatterns"`
	SensitiveHeaders []string `mapstructure:"sensitiveHeaders" json:"sensitiveHeaders" yaml:"sensitiveHeaders"`
}

type RouterDetails struct {
	// 路由地址
	Path string
	// 微服务,驼峰命令属性，需要使用 Tag 标签指定名称
	ServiceId string `mapstructure:"serviceId" json:"serviceId" yaml:"serviceId"`
	// 目标路由
	Url string
	// 是否过滤前缀
	StripPrefix bool `mapstructure:"stripPrefix" json:"stripPrefix" yaml:"stripPrefix"`
	// 超时时间
	Timeout time.Duration
}