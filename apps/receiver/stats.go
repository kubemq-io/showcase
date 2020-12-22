package main

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"go.uber.org/atomic"
	"log"
	"os"
	"time"
)

type StatsInterface interface {
	GetClientStats() *ClientStats
}
type ClientStats struct {
	Messages *atomic.Int64
	Volume   *atomic.Int64
	Errors   *atomic.Int64
}

func NewClientStats() *ClientStats {
	return &ClientStats{
		Messages: atomic.NewInt64(0),
		Volume:   atomic.NewInt64(0),
		Errors:   atomic.NewInt64(0),
	}
}

type Stats struct {
	cfg         *Config
	Source      string
	Group       string
	Instance    string
	Senders     int
	Messages    *atomic.Int64
	Volume      *atomic.Int64
	Errors      *atomic.Int64
	lastSend    int64
	lastReceive int64
	startAt     time.Time
	restyClient *resty.Client
	lastMetric  *Metric
}

func CreateStats(cfg *Config) *Stats {
	hostName, _ := os.Hostname()
	s := &Stats{
		cfg:         cfg,
		Source:      cfg.Source,
		Group:       cfg.Group,
		Instance:    hostName,
		Senders:     cfg.Receivers,
		Messages:    atomic.NewInt64(0),
		Volume:      atomic.NewInt64(0),
		Errors:      atomic.NewInt64(0),
		lastSend:    0,
		lastReceive: 0,
		startAt:     time.Now(),
		restyClient: resty.New(),
		lastMetric: &Metric{
			Source:   cfg.Source,
			Group:    cfg.Group,
			Instance: hostName,
			Clients:  0,
			Messages: 0,
			Volume:   0,
			Errors:   0,
		},
	}
	return s
}

func (s *Stats) CollectStats(clients []StatsInterface) *Stats {
	for _, client := range clients {
		stats := client.GetClientStats()
		s.Messages.Add(stats.Messages.Swap(0))
		s.Volume.Add(stats.Volume.Swap(0))
		s.Errors.Add(stats.Errors.Swap(0))
	}
	return s
}
func (s *Stats) ReportStats() *Stats {
	currentMetric := &Metric{
		Source:   s.Source,
		Group:    s.Group,
		Instance: s.Instance,
		Clients:  s.cfg.Receivers,
		Messages: s.Messages.Load(),
		Volume:   s.Volume.Load(),
		Errors:   s.Errors.Load(),
	}
	reportMetric := &Metric{
		Source:   s.Source,
		Group:    s.Group,
		Instance: s.Instance,
		Clients:  s.cfg.Receivers,
		Messages: currentMetric.Messages - s.lastMetric.Messages,
		Volume:   currentMetric.Volume - s.lastMetric.Volume,
		Errors:   currentMetric.Errors - s.lastMetric.Errors,
	}
	resp, err := s.restyClient.R().SetBody(reportMetric).Post(fmt.Sprintf("%s/report", s.cfg.CollectorUrl))
	if err != nil {
		log.Println(fmt.Sprintf("error reporting stats, %s", err.Error()))
		return s
	}
	if resp.IsError() {
		log.Println(fmt.Sprintf("error reporting stats, %s", resp.Status()))
	}
	s.lastMetric = currentMetric
	return s
}

func (s *Stats) Print() {
	currentSent := s.Messages.Load()
	s.lastSend = currentSent
	log.Println(fmt.Sprintf(
		"Source: %s, "+
			"Group: %s, "+
			"Instance: %s, "+
			"Duration: %s, "+
			"Messages: %d, "+
			"Volume: %d, "+
			"Errors: %d ",
		s.Source,
		s.Group,
		s.Instance,
		time.Since(s.startAt).Round(time.Second),
		s.Messages.Load(),
		s.Volume.Load(),
		s.Errors.Load()))
}

type Metric struct {
	Source   string `json:"source"`
	Group    string `json:"group"`
	Instance string `json:"instance"`
	Clients  int    `json:"clients"`
	Messages int64  `json:"messages"`
	Volume   int64  `json:"volume"`
	Errors   int64  `json:"errors"`
}
