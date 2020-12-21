package collector

import (
	"context"
	"github.com/kubemq-io/showcase/apps/collector/pkg/types"
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
	val, ok := c.Aggregators.Load(m.Source)
	if ok {
		agg := val.(*types.Aggregator)
		agg.Add(m)
	} else {
		agg := types.NewAggregator(m.Source)
		c.Aggregators.Store(m.Source, agg)
		agg.Add(m)
	}

}
func (c *Collector) Clear(source string) {
	val, ok := c.Aggregators.Load(source)
	if ok {
		agg := val.(*types.Aggregator)
		c.snapshot(source, agg)
		agg.Clear()
	}
}
func (c *Collector) ClearAll() {
	c.Aggregators.Range(func(key, value interface{}) bool {
		source, agg := key.(string), value.(*types.Aggregator)
		c.snapshot(source, agg)
		agg.Clear()
		return true
	})
}
func (c *Collector) Top() []*types.Snapshot {
	var list []*types.Snapshot
	c.Buckets.Range(func(key, value interface{}) bool {
		bucket := value.(*types.Bucket)
		list = append(list, bucket.Top())
		return true
	})
	return list
}
func (c *Collector) Bucket(source string, count int) []*types.Snapshot {
	val, ok := c.Buckets.Load(source)
	if ok {
		return val.(*types.Bucket).List(count)
	}
	return nil
}

func (c *Collector) processSnapshots() {
	c.Aggregators.Range(func(key, value interface{}) bool {
		source, agg := key.(string), value.(*types.Aggregator)
		c.snapshot(source, agg)
		return true
	})
}

func (c *Collector) snapshot(source string, agg *types.Aggregator) {
	val, ok := c.Buckets.Load(source)
	if ok {
		bucket := val.(*types.Bucket)
		bucket.Append(agg.Snapshot())
	} else {
		bucket := types.NewBucket(source)
		c.Buckets.Store(source, bucket)
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
