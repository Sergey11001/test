package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Env           string        `yaml:"env" env:"ENV" env-default:"local"`
	DealTimeout   time.Duration `yaml:"dealTimeout" env:"DEAL_TIMEOUT" env-default:"4s"`
	BroadcastPort string        `yaml:"broadcastPort" env:"BROADCAST_PORT" env-default:"12345"`
	GRPCConfig    GRPCConfig    `yaml:"grpc" env:"GRPC" env-required:"true"`
}

type GRPCConfig struct {
	Port    string        `yaml:"port" env:"GRPC_PORT" env-required:"true"`
	Timeout time.Duration `yaml:"timeout" env:"GRPC_TIMEOUT" env-default:"1h"`
}

func MustLoad() *Config {
	path := parseConfigPath()
	if path == "" {
		panic("config path is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file not found: " + path)
	}

	config := &Config{}
	if err := cleanenv.ReadConfig(path, config); err != nil {
		panic("can't read config: " + err.Error())
	}

	return config
}

func parseConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "config file path")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}
	return res
}
