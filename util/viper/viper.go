package viper

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

type config struct {
	Server struct {
		Name string `yaml:"name"`
	} `yaml:"server"`
}

func New() *viper.Viper {
	cf := viper.NewWithOptions(viper.EnvKeyReplacer(strings.NewReplacer(".", "_")))
	cf.SetConfigFile("env/config")
	cf.SetConfigType("yaml")
	cf.AutomaticEnv()
	if err := cf.ReadInConfig(); err != nil {
		log.Fatalf("error loading configuration: %v", err)
	}
	return cf
}
