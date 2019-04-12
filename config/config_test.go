package config_test

import (
	"github.com/sundogrd/user-grpc/utils/config"
	"github.com/sundogrd/user-grpc/utils/file"
	"github.com/sundogrd/user-grpc/utils/pointer"
	"path/filepath"
	"testing"
)

func initConfig() (*config.Config, error) {
	commonConfDir, _ := file.GetCurrentFilePath()
	commonConfDir = filepath.Join(commonConfDir, "../../../config")
	return config.ReadConfigFromFile(commonConfDir, pointer.String("service.config.json"))
}

func TestConfig(t *testing.T) {
	testConfig, err := initConfig()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("config is %+v", testConfig)
}
