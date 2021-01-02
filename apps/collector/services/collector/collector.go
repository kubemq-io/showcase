package collector

import (
	"context"
	"fmt"
	"github.com/kubemq-io/showcase/apps/collector/pkg/types/base"
	"strings"
	"sync"
	"time"
)

const (
	snapshotInterval = 5 * time.Second
)

type Collector struct {
	Aggregators sync.Map
	Buckets     sync.Map
}

func NewCollector(ctx context.Context) (*Collector, error) {
	c := &Collector{
		Aggregators: sync.Map{},
		Buckets:     sync.Map{},
	}
	go c.run(ctx)
	return c, nil
}

func (c *Collector) Aggregate(m *base.Metric) {
	key := fmt.Sprintf("%s/%s", m.Source, m.Group)
	val, ok := c.Aggregators.Load(key)
	if ok {
		agg := val.(*base.Aggregator)
		agg.Add(m)
	} else {
		agg := base.NewAggregator(m.Source, m.Group)
		c.Aggregators.Store(key, agg)
		agg.Add(m)
	}
}
func (c *Collector) Clear(source, group string) {
	key := fmt.Sprintf("%s/%s", source, group)
	val, ok := c.Aggregators.Load(key)
	if ok {
		agg := val.(*base.Aggregator)
		c.snapshot(group, agg)
		agg.Clear()
	}
}
func (c *Collector) ClearAll() {
	//c.Aggregators.Range(func(key, value interface{}) bool {
	//	_, agg := key.(string), value.(*base.Aggregator)
	//	agg.Clear()
	//	return true
	//})
	c.Aggregators = sync.Map{}
	c.Buckets = sync.Map{}
}
func (c *Collector) Top(group string) []*base.Snapshot {
	var list []*base.Snapshot
	c.Buckets.Range(func(key, value interface{}) bool {
		val, bucket := key.(string), value.(*base.Bucket)
		if group == "" {
			list = append(list, bucket.Top())
		} else {
			if strings.Contains(val, "/"+group) {
				list = append(list, bucket.Top())
			}
		}

		return true
	})
	return list
}
func (c *Collector) Bucket(name string, count int) []*base.Snapshot {
	val, ok := c.Buckets.Load(name)
	if ok {
		return val.(*base.Bucket).List(count)
	}
	return nil
}

func (c *Collector) processSnapshots() {
	c.Aggregators.Range(func(key, value interface{}) bool {
		val, agg := key.(string), value.(*base.Aggregator)
		c.snapshot(val, agg)
		return true
	})
}

func (c *Collector) snapshot(key string, agg *base.Aggregator) {
	val, ok := c.Buckets.Load(key)
	if ok {
		bucket := val.(*base.Bucket)
		bucket.Append(agg.Snapshot())
	} else {
		bucket := base.NewBucket(key)
		c.Buckets.Store(key, bucket)
		bucket.Append(agg.Snapshot())
	}
}

func (c *Collector) run(ctx context.Context) {
	for {
		select {
		case <-time.After(snapshotInterval):
			c.processSnapshots()
		case <-ctx.Done():
			return
		}
	}
}
