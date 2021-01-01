package base

import (
	"fmt"
	"github.com/kubemq-io/showcase/apps/collector/pkg/utils"
	"time"
)

type State struct {
	Time      time.Time `json:"time"`
	Instances int       `json:"instances"`
	Clients   int       `json:"clients"`
	Messages  int64     `json:"messages"`
	Volume    int64     `json:"volume"`
	Errors    int64     `json:"errors"`
	Pending   int64     `json:"pending"`
}
type Delta struct {
	Duration  float64 `json:"duration"`
	Instances int     `json:"instances"`
	Clients   int     `json:"clients"`
	Messages  int64   `json:"messages"`
	Volume    int64   `json:"volume"`
	Errors    int64   `json:"errors"`
	Pending   int64   `json:"pending"`
}
type Snapshot struct {
	Source  string `json:"source"`
	Group   string `json:"group"`
	Start   *State `json:"start"`
	End     *State `json:"end"`
	Changed *Delta `json:"changed"`
}

func (s *Snapshot) String() string {
	return fmt.Sprintf(
		"Source: %s, Group: %s, Instances: %d, Clients: %d, Messages: %d, Volume: %s, Errors: %d, Pending: %d",
		s.Source,
		s.Group,
		s.End.Instances,
		s.End.Clients,
		s.End.Messages,
		utils.ByteCount(s.End.Volume),
		s.End.Errors,
		s.End.Pending,
	)
}
func (s *Snapshot) SetDelta() {
	s.Changed = &Delta{
		Duration:  s.End.Time.Sub(s.Start.Time).Seconds(),
		Instances: s.End.Instances - s.Start.Instances,
		Clients:   s.End.Clients - s.Start.Clients,
		Messages:  s.End.Messages - s.Start.Messages,
		Volume:    s.End.Volume - s.Start.Volume,
		Errors:    s.End.Errors - s.Start.Errors,
		Pending:   s.End.Pending - s.Start.Pending,
	}
}
