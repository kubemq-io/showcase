package main

import (
	"context"
	"fmt"
	"github.com/kubemq-io/kubemq-go"
	"github.com/nats-io/nuid"
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
	payload       []byte
	sendChannel   chan []*kubemq.EventStore
}

func NewStoreClient(ctx context.Context, id int, cfg *Config, payload []byte) *StoreClient {
	c := &StoreClient{
		Id:            id,
		cfg:           cfg,
		stats:         NewClientStats(),
		localClientID: fmt.Sprintf("%s-%d", cfg.ClientId, id),
		localChannel:  fmt.Sprintf("%s.%d", cfg.Channel, id),
		sendChannel:   make(chan []*kubemq.EventStore, 60),
		payload:       payload,
	}
	hosts := strings.Split(cfg.Hosts, ",")
	for _, host := range hosts {
		go c.runWorker(ctx, host)
	}
	go c.runGenerator(ctx)

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
	sendCh := make(chan *kubemq.EventStore, c.cfg.SendBatch)
	receiveCh := make(chan *kubemq.EventStoreResult, c.cfg.SendBatch)
	errCh := make(chan error, 1)
	isStreamUp := false

	for {
		if !isStreamUp {
			go client.StreamEventsStore(ctx, sendCh, receiveCh, errCh)
			isStreamUp = true
		}
		quitCh := make(chan bool, 1)
		go func() {
			for {
				select {
				case events := <-c.sendChannel:
					for _, event := range events {
						select {
						case sendCh <- event:
						case <-ctx.Done():
							return
						default:
							c.stats.Errors.Inc()
						}
					}
				case <-quitCh:
					return
				case <-ctx.Done():
					return
				}
			}
		}()
		for {
			select {
			case result := <-receiveCh:
				if result.Sent {
					c.stats.Messages.Inc()
				} else {
					c.Log(result.Err.Error())
					c.stats.Errors.Inc()
				}
			case err := <-errCh:
				c.Logf("sender: %d, error: %s", c.Id, err.Error())
				isStreamUp = false
				goto done
			case <-ctx.Done():
				return
			}
		}
	done:
		time.Sleep(time.Second)
		c.Logf("sender: %d reconnecting", c.Id)
	}

}
func (c *StoreClient) runGenerator(ctx context.Context) {
	for {
		select {
		case <-time.After(time.Duration(c.cfg.SendInterval) * time.Second):
			var events []*kubemq.EventStore
			for j := 0; j < c.cfg.SendBatch; j++ {
				events = append(events,
					kubemq.NewEventStore().
						SetId(nuid.Next()).
						SetChannel(c.localChannel).
						SetBody(c.payload))
			}
			c.sendChannel <- events
		case <-ctx.Done():
			return
		}
	}

}

func (c *StoreClient) GetClientStats() *ClientStats {
	return c.stats
}
