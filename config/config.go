package config

import (
	"github.com/prometheus/common/log"
	"path/filepath"

	"github.com/Jeffail/gabs"
)

type Config struct {
	JsonContainer gabs.Container
}

func (conf Config) Get(path string) interface{} {
	return conf.JsonContainer.Path(path).Data()
}

// NewConfig ...
func NewConfig() Config {
	return Config{}
}

func ReadConfigFromFile(confDir string, configFile *string) (*Config, error) {
	Conf := &Config{}

	if configFile != nil {
		container, err := gabs.ParseJSONFile(filepath.Join(confDir, *configFile))
		if err != nil {
			return nil, err
		}
		Conf.JsonContainer = *container
	} else {
		log.Infof("[config] load from default file service.config.json")
		container, err := gabs.ParseJSONFile(filepath.Join(confDir, "service.config.json"))
		if err != nil {
			return nil, err
		}
		Conf.JsonContainer = *container
	}

	log.Infof("[config] Init config in %s, Conf: %+v\n", confDir, Conf)
	return Conf, nil
}