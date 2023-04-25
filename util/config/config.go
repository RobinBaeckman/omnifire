package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Runtime  RuntimeConfig  `yaml:"runtime"`
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Trace    TraceConfig    `yaml:"trace"`
	Profile  ProfileConfig  `yaml:"profile"`
	NextHop  NextHopConfig  `yaml:"nexthop"`
	TLS      TLSConfig      `yaml:"tls"`
}

type RuntimeConfig struct {
	Env string `yaml:"env"`
}

type ServerConfig struct {
	Name     string `yaml:"name"`
	HttpPort string `yaml:"httpPort"`
	GrpcPort string `yaml:"grpcPort"`
}

type DatabaseConfig struct {
	User          string `yaml:"user"`
	Password      string `yaml:"password"`
	DBName        string `yaml:"dbname"`
	Host          string `yaml:"host"`
	SslMode       string `yaml:"sslmode"`
	MigrationPath string `yaml:"migrationPath"`
}

type TraceConfig struct {
	CollectorHost string `yaml:"collectorHost"`
}

type ProfileConfig struct {
	Host string `yaml:"host"`
}

type NextHopConfig struct {
	Host string `yaml:"host"`
}

type TLSConfig struct {
	Enabled   bool   `yaml:"enabled"`
	ServerCrt string `yaml:"serverCrt"`
	ServerKey string `yaml:"serverKey"`
	ClientCrt string `yaml:"clientCrt"`
	ClientKey string `yaml:"clientKey"`
	CaCrt     string `yaml:"caCrt"`
}

func New() *Config {
	v := viper.NewWithOptions(viper.EnvKeyReplacer(strings.NewReplacer(".", "_")))
	v.SetConfigFile("env/config")
	v.SetConfigType("yaml")
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("error loading configuration: %v", err)
	}
	var cf Config
	if err := v.Unmarshal(&cf); err != nil {
		log.Fatalf("failed to unmarshal config: %v", err)
	}
	return &cf
}
