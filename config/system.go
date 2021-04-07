package config

type System struct {
	DbType        string `mapstructure:"db-type" json:"dbType" yaml:"db-type"`
	UseMultipoint bool   `mapstructure:"use-multipoint" json:"useMultipoint" yaml:"use-multipoint"`
}
