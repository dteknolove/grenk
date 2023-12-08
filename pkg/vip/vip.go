package vip

import (
	"errors"
	"os"

	"github.com/spf13/viper"
)

const (
	ErrReadConfig = "error reading config: "
	StatusProd    = "production"
	StatusDev     = "development"
)

const (
	C_Database = "database" + "."
)

type (
	VipConf struct{}
	VipRes  struct {
		DbName     string
		DbPassword string
		DbUsername string
		DbPort     int
		DbHost     string
		DbSchema   string
		RepoPath   string
	}
)

func New() *VipConf {
	return &VipConf{}
}

func (v *VipConf) config() (*viper.Viper, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, errors.New("error reading workdir: " + err.Error())
	}

	config := viper.New()
	config.SetConfigFile("grenk.yaml")
	config.AddConfigPath(wd)
	errRead := config.ReadInConfig()
	if errRead != nil {
		return nil, errors.New("error reading yaml: " + errRead.Error())
	}
	return config, nil
}

func (v *VipConf) App() (*VipRes, error) {
	config, errConf := v.config()
	if errConf != nil {
		return nil, errConf
	}

	res := &VipRes{
		DbName:     config.GetString(C_Database + "name"),
		DbPassword: config.GetString(C_Database + "password"),
		DbUsername: config.GetString(C_Database + "username"),
		DbPort:     config.GetInt(C_Database + "port"),
		DbHost:     config.GetString(C_Database + "host"),
		DbSchema:   config.GetString(C_Database + "schema"),
		RepoPath:   config.GetString(C_Database + "repo_path"),
	}
	return res, nil
}
