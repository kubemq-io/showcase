package collector

import (
	"context"
	"fmt"
	"github.com/kubemq-io/showcase/apps/collector/pkg/types"
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

func (c *Collector) Aggregate(m *types.Metric) {
	key := fmt.Sprintf("%s/%s", m.Source, m.Group)
	val, ok := c.Aggregators.Load(key)
	if ok {
		agg := val.(*types.Aggregator)
		agg.Add(m)
	} else {
		agg := types.NewAggregator(m.Source, m.Group)
		c.Aggregators.Store(key, agg)
		agg.Add(m)
	}

}
func (c *Collector) Clear(source, group string) {
	key := fmt.Sprintf("%s/%s", source, group)
	val, ok := c.Aggregators.Load(key)
	if ok {
		agg := val.(*types.Aggregator)
		c.snapshot(group, agg)
		agg.Clear()
	}
}
func (c *Collector) ClearAll() {
	c.Aggregators.Range(func(key, value interface{}) bool {
		val, agg := key.(string), value.(*types.Aggregator)
		sourceGroup := strings.Split(val, "/")
		if len(sourceGroup) == 2 {
			c.snapshot(sourceGroup[1], agg)
			agg.Clear()
		}
		return true
	})
}
func (c *Collector) Top(group string) []*types.Snapshot {
	var list []*types.Snapshot
	c.Buckets.Range(func(key, value interface{}) bool {
		groupVal, bucket := key.(string), value.(*types.Bucket)
		if group != "" {
			if group == groupVal {
				list = append(list, bucket.Top())
			}
		} else {
			list = append(list, bucket.Top())
		}
		return true
	})
	return list
}
func (c *Collector) Bucket(name string, count int) []*types.Snapshot {
	val, ok := c.Buckets.Load(name)
	if ok {
		return val.(*types.Bucket).List(count)
	}
	return nil
}

func (c *Collector) processSnapshots() {
	c.Aggregators.Range(func(key, value interface{}) bool {
		val, agg := key.(string), value.(*types.Aggregator)
		sourceGroup := strings.Split(val, "/")
		if len(sourceGroup) == 2 {
			c.snapshot(sourceGroup[1], agg)
		}
		return true
	})
}

func (c *Collector) snapshot(key string, agg *types.Aggregator) {
	val, ok := c.Buckets.Load(key)
	if ok {
		bucket := val.(*types.Bucket)
		bucket.Append(agg.Snapshot())
	} else {
		bucket := types.NewBucket(key)
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
