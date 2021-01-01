package api

import "github.com/kubemq-io/showcase/apps/collector/pkg/types/base"

type Receivers struct {
	Total []*DataItem   `json:"total"`
	List  []*MetricItem `json:"list"`
}

func GetReceivers(data []*base.Snapshot) *Receivers {
	receivers := &Receivers{
		Total: nil,
		List:  []*MetricItem{},
	}
	dataItemsMap := NewDataItemsMap([]string{
		"Clients",
		"Messages",
		"Volume",
		"Errors",
	})
	for _, snapshot := range data {
		mi := NewMetricItem(snapshot.Source).
			SetClients(snapshot.End.Clients).
			SetMessages(snapshot.End.Messages).
			SetVolume(snapshot.End.Volume).
			SetErrors(snapshot.End.Errors)
		receivers.List = append(receivers.List, mi)
		dataItemsMap.Item("Clients").AddIntValue(snapshot.End.Clients)
		dataItemsMap.Item("Messages").AddInt64Value(snapshot.End.Messages)
		dataItemsMap.Item("Volume").AddInt64Value(snapshot.End.Volume)
		dataItemsMap.Item("Errors").AddInt64Value(snapshot.End.Errors)
	}
	dataItemsMap.m["Clients"] = dataItemsMap.Item("Clients").Number()
	dataItemsMap.m["Messages"] = dataItemsMap.Item("Messages").Number()
	dataItemsMap.m["Volume"] = dataItemsMap.Item("Volume").Bytes()
	dataItemsMap.m["Errors"] = dataItemsMap.Item("Errors").Number()
	receivers.Total = dataItemsMap.List()
	return receivers
}
