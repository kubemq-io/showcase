package api

import (
	"github.com/kubemq-io/showcase/apps/collector/pkg/types/kubemq"
	"sort"
	"strings"
)

type Kubemq struct {
	TotalIn  []*DataItem   `json:"total_in"`
	TotalOut []*DataItem   `json:"total_out"`
	List     []*MetricItem `json:"list"`
}

func GetKubeMQ(servers map[string]*kubemq.Status) *Kubemq {
	kubemqs := &Kubemq{
		TotalIn:  nil,
		TotalOut: nil,
		List:     []*MetricItem{},
	}
	dataItemsMapIn := NewDataItemsMap([]string{
		"Channels",
		"Messages",
		"Volume",
		"Errors",
	})
	dataItemsMapOut := NewDataItemsMap([]string{
		"Channels",
		"Messages",
		"Volume",
		"Pending",
		"Errors",
	})
	for _, status := range servers {
		channels := 0
		inMessages := int64(0)
		inVolume := int64(0)
		inErrors := int64(0)
		outMessages := int64(0)
		outVolume := int64(0)
		outErrors := int64(0)
		outPending := int64(0)

		for _, group := range status.Entities {
			channels = group.Total
			inMessages += group.In.Messages
			inVolume += group.In.Volume
			inErrors += group.In.Errors
			outMessages += group.Out.Messages
			outVolume += group.Out.Volume
			outErrors += group.Out.Errors
			outPending += group.Out.Waiting
		}
		mi := NewMetricItem(status.Host).
			SetCPU(status.System.TotalCPUs).
			SetMemory(status.System.ProcessMemory).
			SetMemoryUtilization(status.System.MemoryUtilization).
			SetCPUUtilization(status.System.CPUUtilization).
			SetErrors(inErrors + outErrors).
			SetMessages(inMessages + outMessages).
			SetVolume(inVolume + outVolume).
			SetPending(outPending).
			SetChannels(channels)
		kubemqs.List = append(kubemqs.List, mi)
		dataItemsMapIn.Item("Channels").AddIntValue(channels)
		dataItemsMapIn.Item("Messages").AddInt64Value(inMessages)
		dataItemsMapIn.Item("Volume").AddInt64Value(inVolume)
		dataItemsMapIn.Item("Errors").AddInt64Value(inErrors)

		dataItemsMapOut.Item("Channels").AddIntValue(channels)
		dataItemsMapOut.Item("Messages").AddInt64Value(outMessages)
		dataItemsMapOut.Item("Volume").AddInt64Value(outVolume)
		dataItemsMapOut.Item("Pending").AddInt64Value(outPending)
		dataItemsMapOut.Item("Errors").AddInt64Value(outErrors)

	}

	dataItemsMapIn.m["Channels"] = dataItemsMapIn.Item("Channels").Number()
	dataItemsMapIn.m["Messages"] = dataItemsMapIn.Item("Messages").Number()
	dataItemsMapIn.m["Volume"] = dataItemsMapIn.Item("Volume").Bytes()
	dataItemsMapIn.m["Errors"] = dataItemsMapIn.Item("Errors").Number()

	dataItemsMapOut.m["Channels"] = dataItemsMapOut.Item("Channels").Number()
	dataItemsMapOut.m["Messages"] = dataItemsMapOut.Item("Messages").Number()
	dataItemsMapOut.m["Volume"] = dataItemsMapOut.Item("Volume").Bytes()
	dataItemsMapOut.m["Pending"] = dataItemsMapOut.Item("Pending").Number()
	dataItemsMapOut.m["Errors"] = dataItemsMapOut.Item("Errors").Number()

	kubemqs.TotalIn = dataItemsMapIn.List()
	kubemqs.TotalOut = dataItemsMapOut.List()
	sort.Slice(kubemqs.List, func(i, j int) bool {
		n := strings.Compare(kubemqs.List[i].Title, kubemqs.List[j].Title)
		return n < 0
	})
	return kubemqs
}
