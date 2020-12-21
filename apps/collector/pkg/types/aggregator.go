package types

import (
	"fmt"
	"github.com/kubemq-io/showcase/apps/collector/pkg/logger"
	"go.uber.org/atomic"
	"sync"
	"time"
)

type Aggregator struct {
	sync.Mutex
	logger       *logger.Logger
	LastUpdate   time.Time
	LastSnapshot *Snapshot
	Source       string
	Instances    map[string]int
	Messages     *atomic.Int64
	Volume       *atomic.Int64
	Errors       *atomic.Int64
}

func NewAggregator(source string) *Aggregator {
	return &Aggregator{
		Mutex:  sync.Mutex{},
		logger: logger.NewLogger(fmt.Sprintf("aggregator-%s", source)),
		Source: source,
		LastSnapshot: &Snapshot{
			Source: source,
			Start: &State{
				Time:      time.Now().UTC(),
				Instances: 0,
				Clients:   0,
				Messages:  0,
				Volume:    0,
				Errors:    0,
			},
			End: &State{
				Time:      time.Now().UTC(),
				Instances: 0,
				Clients:   0,
				Messages:  0,
				Volume:    0,
				Errors:    0,
			},
			Changed: &Delta{
				Duration:  0,
				Instances: 0,
				Clients:   0,
				Messages:  0,
				Volume:    0,
				Errors:    0,
			},
		},
		LastUpdate: time.Now().UTC(),
		Instances:  map[string]int{},
		Messages:   atomic.NewInt64(0),
		Volume:     atomic.NewInt64(0),
		Errors:     atomic.NewInt64(0),
	}
}

func (a *Aggregator) Add(m *Metric) {
	a.Lock()
	_, ok := a.Instances[m.Instance]
	if !ok {
		a.Instances[m.Instance] = m.Clients
	}
	a.Unlock()
	a.Volume.Add(m.Volume)
	a.Messages.Add(m.Messages)
	a.Errors.Add(m.Errors)
	a.LastUpdate = time.Now().UTC()
}

func (a *Aggregator) Clear() {
	a.Lock()
	a.Instances = map[string]int{}
	a.Unlock()
	a.Volume.Store(0)
	a.Messages.Store(0)
	a.Errors.Store(0)
	a.logger.Info("aggregator cleared")
}

func (a *Aggregator) Snapshot() *Snapshot {
	a.Lock()
	defer a.Unlock()
	totalClients := 0
	for _, i := range a.Instances {
		totalClients += i
	}
	s := &Snapshot{
		Source: a.Source,
		Start: &State{
			Time:      a.LastSnapshot.End.Time,
			Instances: a.LastSnapshot.End.Instances,
			Clients:   a.LastSnapshot.End.Clients,
			Messages:  a.LastSnapshot.End.Messages,
			Volume:    a.LastSnapshot.End.Volume,
			Errors:    a.LastSnapshot.End.Errors,
		},
		End: &State{
			Time:      time.Now().UTC(),
			Instances: len(a.Instances),
			Clients:   totalClients,
			Messages:  a.Messages.Load(),
			Volume:    a.Volume.Load(),
			Errors:    a.Errors.Load(),
		},
	}
	s.SetDelta()
	a.LastSnapshot = s
	return s
}
