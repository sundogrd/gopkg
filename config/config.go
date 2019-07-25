package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Init ...
func Load(configPath string, configName string) (*viper.Viper, error) {
	config := viper.New()
	config.SetConfigName(configName)
	config.AddConfigPath(configPath)
	err := config.ReadInConfig()
	if err != nil {
		logrus.Errorf("[config.Load] load config failed, err: %s", err.Error())
		return nil, err
	}
	return config, nil
}