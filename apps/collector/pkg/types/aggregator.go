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
	Group        string
	Instances    map[string]int
	Messages     *atomic.Int64
	Volume       *atomic.Int64
	Errors       *atomic.Int64
	Pending      *atomic.Int64
	PendingMap   map[string]*int64
}

func NewAggregator(source, group string) *Aggregator {
	return &Aggregator{
		Mutex:  sync.Mutex{},
		logger: logger.NewLogger(fmt.Sprintf("aggregator-%s-%s", source, group)),
		Source: source,
		Group:  group,
		LastSnapshot: &Snapshot{
			Source: source,
			Group:  group,
			Start: &State{
				Time:      time.Now().UTC(),
				Instances: 0,
				Clients:   0,
				Messages:  0,
				Volume:    0,
				Errors:    0,
				Pending:   0,
			},
			End: &State{
				Time:      time.Now().UTC(),
				Instances: 0,
				Clients:   0,
				Messages:  0,
				Volume:    0,
				Errors:    0,
				Pending:   0,
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
		Pending:    atomic.NewInt64(0),
		PendingMap: map[string]*int64{},
	}
}

func (a *Aggregator) Add(m *Metric) {
	a.Lock()
	_, ok := a.Instances[m.Instance]
	if !ok {
		a.Instances[m.Instance] = m.Clients
	}
	p, ok2 := a.PendingMap[m.Instance]
	if !ok2 {
		p := new(int64)
		*p = m.Pending
		a.PendingMap[m.Instance] = p
	} else {
		*p = m.Pending
	}
	a.Pending.Store(0)
	for _, p := range a.PendingMap {
		a.Pending.Add(*p)
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
	a.PendingMap = map[string]*int64{}
	a.Unlock()
	a.Volume.Store(0)
	a.Messages.Store(0)
	a.Errors.Store(0)
	a.Pending.Store(0)
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
		Group:  a.Group,
		Start: &State{
			Time:      a.LastSnapshot.End.Time,
			Instances: a.LastSnapshot.End.Instances,
			Clients:   a.LastSnapshot.End.Clients,
			Messages:  a.LastSnapshot.End.Messages,
			Volume:    a.LastSnapshot.End.Volume,
			Errors:    a.LastSnapshot.End.Errors,
			Pending:   a.LastSnapshot.End.Pending,
		},
		End: &State{
			Time:      time.Now().UTC(),
			Instances: len(a.Instances),
			Clients:   totalClients,
			Messages:  a.Messages.Load(),
			Volume:    a.Volume.Load(),
			Errors:    a.Errors.Load(),
			Pending:   a.Pending.Load(),
		},
	}
	s.SetDelta()
	a.LastSnapshot = s
	return s
}
