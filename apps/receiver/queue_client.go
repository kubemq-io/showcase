package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

type QueueClient struct {
	Id            int
	cfg           *Config
	stats         *ClientStats
	localClientID string
	localChannel  string
}

func NewQueueClient(ctx context.Context, id int, cfg *Config) *QueueClient {
	c := &QueueClient{
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
func (c *QueueClient) Log(msg string) {
	if c.cfg.Verbose {
		log.Println(msg)
	}
}
func (c *QueueClient) Logf(format string, args ...interface{}) {
	if c.cfg.Verbose {
		log.Println(fmt.Sprintf(format, args...))
	}
}
func (c *QueueClient) runWorker(ctx context.Context, address string) {
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
	for {
		results, err := client.NewReceiveQueueMessagesRequest().
			SetChannel(c.localChannel).
			SetMaxNumberOfMessages(c.cfg.ReceiveBatch).
			SetWaitTimeSeconds(c.cfg.ReceiveTimeout).
			Send(ctx)
		if err != nil {
			c.Logf("instance:%d,  error: %s", c.Id, err.Error())
			time.Sleep(time.Second)
			c.stats.Errors.Add(int64(c.cfg.ReceiveBatch))
			continue
		}
		if results != nil {
			if results.IsError {
				c.stats.Errors.Add(int64(c.cfg.ReceiveBatch))
			} else {
				totalVol := 0
				for _, msg := range results.Messages {
					totalVol += len(msg.Body)
					msg = nil
				}
				c.stats.Messages.Add(int64(len(results.Messages)))
				c.stats.Volume.Add(int64(totalVol))
			}

		}
		select {
		case <-ctx.Done():
			return
		default:

		}
	}

}

func (c *QueueClient) GetClientStats() *ClientStats {
	return c.stats
}
