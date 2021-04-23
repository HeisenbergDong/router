package config

type Grpc struct {
	JwtAddress string `mapstructure:"jwt-address" json:"jwtAddress" yaml:"jwt-address"`
}
