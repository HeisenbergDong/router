package config

import "time"

type Server struct {
	Port        int           `mapstructure:"port" json:"port" yaml:"port"`
	ContextPath string        `mapstructure:"contextPath" json:"contextPath" yaml:"contextPath"`
	Timeout     time.Duration `mapstructure:"timeout" json:"timeout" yaml:"timeout"`
	Host        string        `mapstructure:"host" json:"v" yaml:"host"`
}