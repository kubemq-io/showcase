package console

import (
	"context"
	"github.com/kubemq-io/showcase/apps/collector/services/collector"
	"log"
	"time"
)

type Console struct {
	collector *collector.Collector
}

func NewConsole(ctx context.Context, collector *collector.Collector) *Console {
	c := &Console{
		collector: collector,
	}
	go c.run(ctx)
	return c
}
func (c *Console) print() {
	for _, snapshot := range c.collector.Top("") {
		log.Println(snapshot.String())
	}
}
func (c *Console) run(ctx context.Context) {
	for {
		select {
		case <-time.After(5 * time.Second):
			c.print()
		case <-ctx.Done():
			return
		}
	}
}
