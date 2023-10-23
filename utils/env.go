package utils

import (
	"github.com/spf13/viper"
)

func LoadEnvConfig(prefix string) *viper.Viper {
	v := viper.New()
	v.SetEnvPrefix(prefix)

	v.AutomaticEnv()
	return v
}
