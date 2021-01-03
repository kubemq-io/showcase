package config

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	KubeMQHosts   string
	ApiServerUrl  string
	ApiServerPort int
	Console       bool
}

var (
	_ = pflag.String("kubemq-hosts", "", "set kubemq hosts status collection")
	_ = pflag.String("api-server-url", "http://localhost:8085", "set api endpoint address")
	_ = pflag.Int("api-server-port", 8085, "set api server port")
	_ = pflag.Bool("console", false, "export to console")
)

func LoadConfig() (*Config, error) {
	pflag.Parse()
	cfg := &Config{}
	viper.BindEnv("KubeMQHosts", "KUBEMQ_HOSTS")
	viper.BindEnv("ApiServerUrl", "API_SERVER_URL")
	viper.BindEnv("ApiServerPort", "API_SERVER_PORT")
	viper.BindEnv("Console", "CONSOLE")

	viper.BindPFlag("KubeMQHosts", pflag.CommandLine.Lookup("kubemq-hosts"))
	viper.BindPFlag("ApiServerUrl", pflag.CommandLine.Lookup("api-server-url"))
	viper.BindPFlag("ApiServerPort", pflag.CommandLine.Lookup("api-server-port"))
	viper.BindPFlag("Console", pflag.CommandLine.Lookup("console"))

	err := viper.Unmarshal(cfg)
	if err != nil {
		return nil, err
	}
	wc := &WebConfig{
		ApiServerURL: cfg.ApiServerUrl,
		PollInterval: 5000,
	}
	err = wc.Write("./dist/runtimeConfig.json")
	return cfg, nil
}

func (c *Config) Print() {
	log.Println("KubeMQHosts->", c.KubeMQHosts)
	log.Println("ApiServerUrl->", c.ApiServerUrl)
	log.Println("ApiServerPort->", c.ApiServerPort)
	log.Println("Console->", c.Console)
}
