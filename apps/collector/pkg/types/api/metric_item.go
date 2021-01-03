package api

import (
	"fmt"
	"github.com/dustin/go-humanize"
)

type MetricItem struct {
	Title             string `json:"title"`
	Clients           string `json:"clients,omitempty"`
	Channels          string `json:"channels,omitempty"`
	Messages          string `json:"messages,omitempty"`
	Volume            string `json:"volume,omitempty"`
	Errors            string `json:"errors,omitempty"`
	Pending           string `json:"pending,omitempty"`
	CPU               string `json:"cpu,omitempty"`
	CPUUtilization    string `json:"cpu_utilization,omitempty"`
	Memory            string `json:"memory,omitempty"`
	MemoryUtilization string `json:"memory_utilization,omitempty"`
}

func NewMetricItem(title string) *MetricItem {
	return &MetricItem{
		Title: title,
	}
}

func (m *MetricItem) SetClients(value int) *MetricItem {
	m.Clients = fmt.Sprintf("%d", value)
	return m
}
func (m *MetricItem) SetChannels(value int) *MetricItem {
	m.Channels = fmt.Sprintf("%d", value)
	return m
}

func (m *MetricItem) SetMessages(value int64) *MetricItem {
	m.Messages = humanize.Comma(value)
	return m
}
func (m *MetricItem) SetVolume(value int64) *MetricItem {
	m.Volume = humanize.Bytes(uint64(value))
	return m
}
func (m *MetricItem) SetErrors(value int64) *MetricItem {
	m.Errors = humanize.Comma(value)
	return m
}
func (m *MetricItem) SetPending(value int64) *MetricItem {
	m.Pending = humanize.Comma(value)
	return m
}
func (m *MetricItem) SetCPU(value int) *MetricItem {
	m.CPU = fmt.Sprintf("%d", value)
	return m
}
func (m *MetricItem) SetCPUUtilization(value float64) *MetricItem {
	m.CPUUtilization = fmt.Sprintf("%s", humanize.CommafWithDigits(value, 2))
	return m
}
func (m *MetricItem) SetMemory(value float64) *MetricItem {
	m.Memory = humanize.Bytes(uint64(value))
	return m
}
func (m *MetricItem) SetMemoryUtilization(value float64) *MetricItem {
	m.MemoryUtilization = fmt.Sprintf("%s", humanize.CommafWithDigits(value, 2))
	return m
}
