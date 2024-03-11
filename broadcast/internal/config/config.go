package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	BroadcastIP   string `envconfig:"BROADCAST_IP" required:"true"`
	BroadcastPort int    `envconfig:"BROADCAST_PORT" required:"true"`
	PrefixIP      string `envconfig:"BROADCAST_PREFIX"`
}

func MustLoad() *Config {
	config := &Config{}
	envconfig.MustProcess("", config)

	return config
}
