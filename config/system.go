package config

type System struct {
	DbType        string `mapstructure:"db-type" json:"dbType" yaml:"db-type"`
}
