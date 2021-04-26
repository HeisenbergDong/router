package config

type Config struct {
	System        System   `mapstructure:"system" json:"system" yaml:"system"`
	Server        Server   `mapstructure:"server" json:"server" yaml:"server"`
	GatewayRouter *Routers `mapstructure:"gateway" json:"gateway" yaml:"gateway"`
	Zap           Zap      `mapstructure:"zap" json:"zap" yaml:"zap"`
	Mysql         Mysql    `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	Redis         Redis    `mapstructure:"redis" json:"redis" yaml:"redis"`
	JWT           JWT      `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
	Grpc          Grpc     `mapstructure:"grpc" json:"grpc" yaml:"grpc"`
}
