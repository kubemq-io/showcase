package api

import (
	"github.com/kubemq-io/showcase/apps/collector/pkg/types/base"
	"sort"
	"strings"
)

type Senders struct {
	Total []*DataItem   `json:"total"`
	List  []*MetricItem `json:"list"`
}

func GetSenders(data []*base.Snapshot) *Senders {
	senders := &Senders{
		Total: nil,
		List:  []*MetricItem{},
	}
	dataItemsMap := NewDataItemsMap([]string{
		"Clients",
		"Messages",
		"Volume",
		"Pending",
		"Errors",
	})
	for _, snapshot := range data {
		mi := NewMetricItem(snapshot.Source).
			SetClients(snapshot.End.Clients).
			SetMessages(snapshot.End.Messages).
			SetVolume(snapshot.End.Volume).
			SetErrors(snapshot.End.Errors).
			SetPending(snapshot.End.Pending)

		senders.List = append(senders.List, mi)
		dataItemsMap.Item("Clients").AddIntValue(snapshot.End.Clients)
		dataItemsMap.Item("Messages").AddInt64Value(snapshot.End.Messages)
		dataItemsMap.Item("Volume").AddInt64Value(snapshot.End.Volume)
		dataItemsMap.Item("Pending").AddInt64Value(snapshot.End.Pending)
		dataItemsMap.Item("Errors").AddInt64Value(snapshot.End.Errors)
	}
	dataItemsMap.m["Clients"] = dataItemsMap.Item("Clients").Number()
	dataItemsMap.m["Messages"] = dataItemsMap.Item("Messages").Number()
	dataItemsMap.m["Volume"] = dataItemsMap.Item("Volume").Bytes()
	dataItemsMap.m["Pending"] = dataItemsMap.Item("Pending").Number()
	dataItemsMap.m["Errors"] = dataItemsMap.Item("Errors").Number()
	sort.Slice(senders.List, func(i, j int) bool {
		n := strings.Compare(senders.List[i].Title, senders.List[j].Title)
		return n < 0
	})
	senders.Total = dataItemsMap.List()
	return senders
}
