package config

import (
	"github.com/spf13/viper"
)

type Conf struct {
	RateLimitIP            int    `mapstructure:"RATE_LIMIT_IP"`
	RateLimitToken         int    `mapstructure:"RATE_LIMIT_TOKEN"`
	TimeLimitType          string `mapstructure:"TIME_LIMIT_TYPE"`
	TimeBlockType          string `mapstructure:"TIME_BLOCK_TYPE"`
	BlockLimitTimeDuration int    `mapstructure:"BLOCK_LIMIT_TIME_DURATION"`
}

func LoadConfig(path string) (*Conf, error) {
	var cfg *Conf
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}

	return cfg, err
}
