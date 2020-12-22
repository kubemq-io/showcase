package main

import (
	"context"
	"fmt"
	"github.com/kubemq-io/kubemq-go"
	"log"
	"strconv"
	"strings"

	"time"
)

type StoreClient struct {
	Id            int
	cfg           *Config
	stats         *ClientStats
	localClientID string
	localChannel  string
}

func NewStoreClient(ctx context.Context, id int, cfg *Config) *StoreClient {
	c := &StoreClient{
		Id:            id,
		cfg:           cfg,
		stats:         NewClientStats(),
		localClientID: fmt.Sprintf("%s-%d", cfg.ClientId, id),
		localChannel:  fmt.Sprintf("%s.%d", cfg.Channel, id),
	}
	hosts := strings.Split(cfg.Hosts, ",")
	for _, host := range hosts {
		go c.runWorker(ctx, host)
	}
	return c
}

func (c *StoreClient) Log(msg string) {
	if c.cfg.Verbose {
		log.Println(msg)
	}
}
func (c *StoreClient) Logf(format string, args ...interface{}) {
	if c.cfg.Verbose {
		log.Println(fmt.Sprintf(format, args...))
	}
}

func (c *StoreClient) runWorker(ctx context.Context, address string) {
	host := ""
	port := 0
	parts := strings.Split(address, ":")
	if len(parts) == 2 {
		host = parts[0]
		port, _ = strconv.Atoi(parts[1])
	} else {
		c.Log(fmt.Sprintf("kubemq client bad address, %s", address))
		return
	}

	client, err := getClient(ctx, host, port, c.localClientID)
	if err != nil {
		c.Log(fmt.Sprintf("error get kubemq client: %s", err.Error()))
		return
	}
	defer func() {
		_ = client.Close()
		c.Log("kubemq client complete")
	}()
	errCh := make(chan error, 1)

	for {
		receiveCh, err := client.SubscribeToEventsStore(ctx, c.localChannel, c.cfg.ReceiveGroup, errCh, kubemq.StartFromFirstEvent())
		if err != nil {
			c.Logf("instance:%d, error: %s", c.Id, err.Error())
			c.stats.Errors.Inc()
			time.Sleep(time.Second)
			continue
		}
		for {
			select {
			case msg := <-receiveCh:
				c.stats.Messages.Inc()
				c.stats.Volume.Add(int64(len(msg.Body)))
				msg = nil
			case <-errCh:
				c.stats.Errors.Inc()
			case <-ctx.Done():
				return
			}
		}
	}

}

func (c *StoreClient) GetClientStats() *ClientStats {
	return c.stats
}
