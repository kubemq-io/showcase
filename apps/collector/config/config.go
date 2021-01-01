package config

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	KubeMQHosts string
}

var (
	_ = pflag.String("kubemq-hosts", "localhost:5000", "set kubemq hosts status collection")
)

func LoadConfig() (*Config, error) {
	pflag.Parse()
	cfg := &Config{}
	viper.BindEnv("KubeMQHosts", "KUBEMQ_HOSTS")
	viper.BindPFlag("KubeMQHosts", pflag.CommandLine.Lookup("kubemq-hosts"))
	err := viper.Unmarshal(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) Print() {
	log.Println("KubeMQHosts->", c.KubeMQHosts)

}
