package types

import (
	"fmt"
	"time"
)

type State struct {
	Time      time.Time `json:"time"`
	Instances int       `json:"instances"`
	Clients   int       `json:"clients"`
	Messages  int64     `json:"messages"`
	Volume    int64     `json:"volume"`
	Errors    int64     `json:"errors"`
}
type Delta struct {
	Duration  float64
	Instances int
	Clients   int
	Messages  int64
	Volume    int64
	Errors    int64
}
type Snapshot struct {
	Source  string
	Start   *State
	End     *State
	Changed *Delta
}

func (s *Snapshot) String() string {
	return fmt.Sprintf(
		"Source: %s, Instances: %d, Senders: %d, Messages: %d, Volume: %d, Errors: %d",
		s.Source,
		s.End.Instances,
		s.End.Clients,
		s.End.Messages,
		s.End.Volume,
		s.End.Errors,
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
	}
}
