package config

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"os"
	"time"
)

type Config struct {
	Env             string        `envconfig:"ENV"`
	DealTimeout     time.Duration `envconfig:"DEAL_TIMEOUT" default:"4s"`
	BroadcastPort   string        `envconfig:"BROADCAST_PORT" required:"true"`
	BroadcastPrefix string        `envconfig:"BROADCAST_PREFIX"`
	GRPCConfig      GRPCConfig
	ConfigPrefix    string `envconfig:"CONFIG_PREFIX"`
}

type GRPCConfig struct {
	Port    string
	Timeout time.Duration `envconfig:"GRPC_TIMEOUT" default:"30s"`
}

func MustLoad() *Config {
	config := &Config{}

	envconfig.MustProcess("", config)

	GRPCPort := os.Getenv(fmt.Sprintf(config.ConfigPrefix + "_GRPC_PORT"))
	if GRPCPort == "" {
		GRPCPort = "50055"
	}
	config.GRPCConfig.Port = GRPCPort

	return config
}
