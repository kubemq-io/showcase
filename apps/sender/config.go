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
	Senders           int
	SendBatch         int
	SendInterval      int
	LoadInterval      int
	KillAfter         int
	PayloadSize       int
	PayloadFile       string
	CollectEvery      int
	Verbose           bool
	CollectorUrl      string
}

var (
	_ = pflag.String("source", "sender", "set source")
	_ = pflag.String("group", "senders", "set group")
	_ = pflag.String("hosts", "localhost:50000", "set hosts")
	_ = pflag.String("channel", nuid.Next(), "set channel")
	_ = pflag.String("type", "store", "set loader type")
	_ = pflag.String("clientId", "test-command-client-id", "set clientId")
	_ = pflag.Int("senders", 100, "set senders")
	_ = pflag.Int("channel-start-range", 0, "set channel start range")
	_ = pflag.Int("sendBatch", 1, "set sendBatch")
	_ = pflag.Int("sendInterval", 1, "set sendInterval")
	_ = pflag.Int("loadInterval", 100, "set loadInterval")
	_ = pflag.Int("killAfter", 0, "set killAfter")
	_ = pflag.Int("payloadSize", 100, "set payloadSize")
	_ = pflag.String("payloadFile", "", "set payload file")
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
	viper.BindEnv("Senders", "SENDERS")
	viper.BindEnv("SendBatch", "SEND_BATCH")
	viper.BindEnv("SendInterval", "SEND_INTERVAL")
	viper.BindEnv("LoadInterval", "LOAD_INTERVAL")
	viper.BindEnv("KillAfter", "KILL_AFTER")
	viper.BindEnv("PayloadSize", "PAYLOAD_SIZE")
	viper.BindEnv("PayloadFile", "PAYLOAD_File")
	viper.BindEnv("CollectEvery", "COLLECT_EVERY")
	viper.BindEnv("Verbose", "VERBOSE")
	viper.BindEnv("Verbose", "VERBOSE")
	viper.BindEnv("CollectorUrl", "COLLECTOR-URL")

	viper.BindPFlag("Source", pflag.CommandLine.Lookup("source"))
	viper.BindPFlag("Group", pflag.CommandLine.Lookup("group"))
	viper.BindPFlag("Hosts", pflag.CommandLine.Lookup("hosts"))
	viper.BindPFlag("Type", pflag.CommandLine.Lookup("type"))
	viper.BindPFlag("Channel", pflag.CommandLine.Lookup("channel"))
	viper.BindPFlag("ChannelStartRange", pflag.CommandLine.Lookup("channel-start-range"))
	viper.BindPFlag("ClientId", pflag.CommandLine.Lookup("clientId"))
	viper.BindPFlag("Senders", pflag.CommandLine.Lookup("senders"))
	viper.BindPFlag("SendBatch", pflag.CommandLine.Lookup("sendBatch"))
	viper.BindPFlag("SendInterval", pflag.CommandLine.Lookup("sendInterval"))
	viper.BindPFlag("LoadInterval", pflag.CommandLine.Lookup("loadInterval"))
	viper.BindPFlag("KillAfter", pflag.CommandLine.Lookup("killAfter"))
	viper.BindPFlag("PayloadSize", pflag.CommandLine.Lookup("payloadSize"))
	viper.BindPFlag("PayloadFile", pflag.CommandLine.Lookup("payloadFile"))
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
	log.Println("Senders->", c.Senders)
	log.Println("SendBatch->", c.SendBatch)
	log.Println("SendInterval->", c.SendInterval)
	log.Println("LoadInterval->", c.LoadInterval)
	log.Println("KillAfter->", c.KillAfter)
	log.Println("PayloadSize->", c.PayloadSize)
	log.Println("PayloadFile->", c.PayloadFile)
}
