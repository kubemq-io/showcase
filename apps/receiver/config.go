package main

import (
	"github.com/nats-io/nuid"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	Source            string
	Group             string
	Hosts             string
	Type              string
	Channel           string
	ChannelStartRange int
	ClientId          string
	Receivers         int
	Concurrency       int
	ReceiveBatch      int
	ReceiveTimeout    int
	ReceiveGroup      string
	LoadInterval      int
	KillAfter         int
	CollectEvery      int
	Verbose           bool
	CollectorUrl      string
}

var (
	_ = pflag.String("source", "receiver", "set source")
	_ = pflag.String("group", "receivers", "set group")
	_ = pflag.String("hosts", "localhost:50000", "set hosts")
	_ = pflag.String("channel", nuid.Next(), "set channel")
	_ = pflag.String("type", "store", "set loader type")
	_ = pflag.String("clientId", "test-command-client-id", "set clientId")
	_ = pflag.Int("receivers", 100, "set receivers")
	_ = pflag.Int("concurrency", 1, "set receivers concurrency")
	_ = pflag.Int("channel-start-range", 0, "set channel start range")
	_ = pflag.Int("receiveBatch", 1, "set sendBatch")
	_ = pflag.Int("receiveTimeout", 60, "set receive timeout")
	_ = pflag.String("receiveGroup", "g1", "set receive group")
	_ = pflag.Int("loadInterval", 100, "set loadInterval")
	_ = pflag.Int("killAfter", 0, "set killAfter")
	_ = pflag.Int("collectEvery", 5, "set collectEvery")
	_ = pflag.Bool("verbose", false, "set verbose")
	_ = pflag.String("collector-url", "http://localhost:8085", "set collector url")
)

func LoadConfig() (*Config, error) {
	pflag.Parse()
	cfg := &Config{}
	viper.BindEnv("Source", "SOURCE")
	viper.BindEnv("Group", "GROUP")
	viper.BindEnv("Hosts", "HOSTS")
	viper.BindEnv("Channel", "CHANNEL")
	viper.BindEnv("Type", "TYPE")
	viper.BindEnv("ClientId", "CLIENT_ID")
	viper.BindEnv("ChannelStartRange", "CHANNEL-START-RANGE")
	viper.BindEnv("Receivers", "RECEIVERS")
	viper.BindEnv("Concurrency", "CONCURRENCY")
	viper.BindEnv("ReceiveBatch", "RECEIVE_BATCH")
	viper.BindEnv("ReceiveTimeout", "RECEIVE_TIMEOUT")
	viper.BindEnv("ReceiveGroup", "RECEIVE_GROUP")
	viper.BindEnv("LoadInterval", "LOAD_INTERVAL")
	viper.BindEnv("KillAfter", "KILL_AFTER")
	viper.BindEnv("CollectEvery", "COLLECT_EVERY")
	viper.BindEnv("Verbose", "VERBOSE")
	viper.BindEnv("CollectorUrl", "COLLECTOR-URL")

	viper.BindPFlag("Source", pflag.CommandLine.Lookup("source"))
	viper.BindPFlag("Group", pflag.CommandLine.Lookup("group"))
	viper.BindPFlag("Hosts", pflag.CommandLine.Lookup("hosts"))
	viper.BindPFlag("Type", pflag.CommandLine.Lookup("type"))
	viper.BindPFlag("Channel", pflag.CommandLine.Lookup("channel"))
	viper.BindPFlag("ChannelStartRange", pflag.CommandLine.Lookup("channel-start-range"))
	viper.BindPFlag("ClientId", pflag.CommandLine.Lookup("clientId"))
	viper.BindPFlag("Receivers", pflag.CommandLine.Lookup("receivers"))
	viper.BindPFlag("Concurrency", pflag.CommandLine.Lookup("concurrency"))
	viper.BindPFlag("ReceiveBatch", pflag.CommandLine.Lookup("receiveBatch"))
	viper.BindPFlag("ReceiveTimeout", pflag.CommandLine.Lookup("receiveTimeout"))
	viper.BindPFlag("ReceiveGroup", pflag.CommandLine.Lookup("receiveGroup"))
	viper.BindPFlag("LoadInterval", pflag.CommandLine.Lookup("loadInterval"))
	viper.BindPFlag("KillAfter", pflag.CommandLine.Lookup("killAfter"))
	viper.BindPFlag("CollectEvery", pflag.CommandLine.Lookup("collectEvery"))
	viper.BindPFlag("Verbose", pflag.CommandLine.Lookup("verbose"))
	viper.BindPFlag("CollectorUrl", pflag.CommandLine.Lookup("collector-url"))

	err := viper.Unmarshal(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) Print() {
	log.Println("Source->", c.Source)
	log.Println("Group->", c.Group)
	log.Println("Hosts->", c.Hosts)
	log.Println("Type->", c.Type)
	log.Println("Channel->", c.Channel)
	log.Println("Channel Start Range->", c.ChannelStartRange)
	log.Println("ClientId->", c.ClientId)
	log.Println("Receivers->", c.Receivers)
	log.Println("Concurrency->", c.Concurrency)
	log.Println("ReceiveBatch->", c.ReceiveBatch)
	log.Println("ReceiveTimout->", c.ReceiveTimeout)
	log.Println("ReceiveGroup->", c.ReceiveGroup)
	log.Println("LoadInterval->", c.LoadInterval)
	log.Println("KillAfter->", c.KillAfter)
	log.Println("CollectEvery->", c.CollectEvery)
	log.Println("Verbose->", c.Verbose)
	log.Println("CollectorUrl->", c.CollectorUrl)

}
