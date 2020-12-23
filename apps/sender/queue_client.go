package main

import (
	"context"
	"fmt"
	"github.com/kubemq-io/kubemq-go"
	"github.com/nats-io/nuid"
	"go.uber.org/atomic"
	"log"
	"strconv"
	"strings"
	"time"
)

type QueueClient struct {
	Id                int
	cfg               *Config
	stats             *ClientStats
	localClientID     string
	localChannel      string
	payload           []byte
	totalMessagesLeft *atomic.Int64
	sendChannel       chan []*kubemq.QueueMessage
}

func NewQueueClient(ctx context.Context, id int, cfg *Config, payload []byte) *QueueClient {
	totalMessages := cfg.TotalMessages
	if totalMessages == 0 {
		totalSeconds := cfg.Duration.Round(time.Second).Seconds()
		mps := 1e3 / cfg.SendInterval * cfg.SendBatch
		totalMessages = int64(totalSeconds * float64(mps))
	}

	c := &QueueClient{
		Id:                id,
		cfg:               cfg,
		stats:             NewClientStats(),
		localClientID:     fmt.Sprintf("%s-%d", cfg.ClientId, id),
		localChannel:      fmt.Sprintf("%s.%d", cfg.Channel, id),
		sendChannel:       make(chan []*kubemq.QueueMessage, 60),
		payload:           payload,
		totalMessagesLeft: atomic.NewInt64(totalMessages),
	}
	hosts := strings.Split(cfg.Hosts, ",")
	for _, host := range hosts {
		for i := 0; i < cfg.Concurrency; i++ {
			go c.runWorker(ctx, host)
		}
	}
	go c.runGenerator(ctx)
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
		select {
		case messages := <-c.sendChannel:
			batch := client.NewQueueMessages()
			batch.Messages = messages
			results, err := batch.Send(ctx)
			if err != nil {
				c.stats.Errors.Add(int64(len(messages)))
				c.totalMessagesLeft.Add(int64(len(messages)))
				c.Logf("error sending queue messages, %s", err.Error())
			} else {
				for _, result := range results {
					if result.IsError {
						c.Logf("error sending queue messages, %s", result.Error)
						c.stats.Errors.Inc()
						c.totalMessagesLeft.Inc()
					} else {
						c.stats.Messages.Inc()
						c.stats.Volume.Add(int64(len(c.payload)))
					}
				}
			}
			c.stats.Pending.Store(c.totalMessagesLeft.Load())
		case <-ctx.Done():
			return
		}
	}

}
func (c *QueueClient) runGenerator(ctx context.Context) {
	for {
		select {
		case <-time.After(time.Duration(c.cfg.SendInterval) * time.Millisecond):
			left := c.totalMessagesLeft.Load()
			if left > 0 {
				wanted := c.cfg.SendBatch
				if left < int64(wanted) {
					wanted = int(left)
				}
				var queueMessages []*kubemq.QueueMessage
				for j := 0; j < wanted; j++ {
					queueMessages = append(queueMessages,
						kubemq.NewQueueMessage().
							SetId(nuid.Next()).
							SetChannel(c.localChannel).
							SetBody(c.payload))
				}
				c.totalMessagesLeft.Sub(int64(wanted))
				c.sendChannel <- queueMessages
			}
		case <-ctx.Done():
			return
		}
	}

}

func (c *QueueClient) GetClientStats() *ClientStats {
	return c.stats
}
