package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

type Config struct {
	BroadcastIP   string `yaml:"broadcastIP" env:"BROADCAST_IP" env-required:"true"`
	BroadcastPort int    `yaml:"broadcastPort" env:"BROADCAST_PORT" env-required:"true"`
	PrefixIP      string `yaml:"prefixIP" env:"PREFIX_IP"`
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
