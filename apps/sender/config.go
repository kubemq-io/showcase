package main

import (
	"github.com/nats-io/nuid"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
	"time"
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
	Concurrency       int
	SendBatch         int
	SendInterval      int
	LoadInterval      int
	Duration          time.Duration
	TotalMessages     int64
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
	_ = pflag.Int("concurrency", 1, "set senders concurrency")
	_ = pflag.Int("channel-start-range", 0, "set channel start range")
	_ = pflag.Int("sendBatch", 1, "set sendBatch")
	_ = pflag.Int64("totalMessages", 0, "set totalMessages")
	_ = pflag.Int("sendInterval", 1000, "set sendInterval")
	_ = pflag.Int("loadInterval", 100, "set loadInterval")
	_ = pflag.Int("killAfter", 0, "set killAfter")
	_ = pflag.Int("payloadSize", 100, "set payloadSize")
	_ = pflag.String("payloadFile", "", "set payload file")
	_ = pflag.Int("collectEvery", 5, "set collectEvery")
	_ = pflag.Bool("verbose", false, "set verbose")
	_ = pflag.String("collector-url", "http://localhost:8085", "set collector url")
	_ = pflag.Duration("duration", time.Hour, "set sending duration")
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
	viper.BindEnv("Concurrency", "CONCURRENCY")
	viper.BindEnv("SendBatch", "SEND_BATCH")
	viper.BindEnv("TotalMessages", "TOTAL_MESSAGES")
	viper.BindEnv("SendInterval", "SEND_INTERVAL")
	viper.BindEnv("LoadInterval", "LOAD_INTERVAL")
	viper.BindEnv("KillAfter", "KILL_AFTER")
	viper.BindEnv("PayloadSize", "PAYLOAD_SIZE")
	viper.BindEnv("PayloadFile", "PAYLOAD_File")
	viper.BindEnv("CollectEvery", "COLLECT_EVERY")
	viper.BindEnv("Verbose", "VERBOSE")
	viper.BindEnv("CollectorUrl", "COLLECTOR-URL")
	viper.BindEnv("Duration", "DURATION")

	viper.BindPFlag("Source", pflag.CommandLine.Lookup("source"))
	viper.BindPFlag("Group", pflag.CommandLine.Lookup("group"))
	viper.BindPFlag("Hosts", pflag.CommandLine.Lookup("hosts"))
	viper.BindPFlag("Type", pflag.CommandLine.Lookup("type"))
	viper.BindPFlag("Channel", pflag.CommandLine.Lookup("channel"))
	viper.BindPFlag("ChannelStartRange", pflag.CommandLine.Lookup("channel-start-range"))
	viper.BindPFlag("ClientId", pflag.CommandLine.Lookup("clientId"))
	viper.BindPFlag("Senders", pflag.CommandLine.Lookup("senders"))
	viper.BindPFlag("Concurrency", pflag.CommandLine.Lookup("concurrency"))
	viper.BindPFlag("TotalMessages", pflag.CommandLine.Lookup("totalMessages"))
	viper.BindPFlag("SendBatch", pflag.CommandLine.Lookup("sendBatch"))
	viper.BindPFlag("SendInterval", pflag.CommandLine.Lookup("sendInterval"))
	viper.BindPFlag("LoadInterval", pflag.CommandLine.Lookup("loadInterval"))
	viper.BindPFlag("KillAfter", pflag.CommandLine.Lookup("killAfter"))
	viper.BindPFlag("PayloadSize", pflag.CommandLine.Lookup("payloadSize"))
	viper.BindPFlag("PayloadFile", pflag.CommandLine.Lookup("payloadFile"))
	viper.BindPFlag("CollectEvery", pflag.CommandLine.Lookup("collectEvery"))
	viper.BindPFlag("Verbose", pflag.CommandLine.Lookup("verbose"))
	viper.BindPFlag("CollectorUrl", pflag.CommandLine.Lookup("collector-url"))
	viper.BindPFlag("Duration", pflag.CommandLine.Lookup("duration"))

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
	log.Println("Concurrency->", c.Concurrency)
	log.Println("SendBatch->", c.SendBatch)
	log.Println("SendInterval->", c.SendInterval)
	log.Println("Duration->", c.Duration)
	log.Println("TotalMessages->", c.TotalMessages)
	log.Println("LoadInterval->", c.LoadInterval)
	log.Println("KillAfter->", c.KillAfter)
	log.Println("PayloadSize->", c.PayloadSize)
	log.Println("PayloadFile->", c.PayloadFile)
	log.Println("CollectEvery->", c.CollectEvery)
	log.Println("Verbose->", c.Verbose)
	log.Println("CollectorUrl->", c.CollectorUrl)
}
