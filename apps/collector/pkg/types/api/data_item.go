package api

import (
	"fmt"
	"github.com/dustin/go-humanize"
)

type DataItem struct {
	Key          string `json:"key"`
	Value        string `json:"value"`
	int64Value   *int64
	float64Value *float64
	intValue     *int
}

func NewDataItem(key string) *DataItem {
	return &DataItem{
		Key:   key,
		Value: "0",
	}
}

func (d *DataItem) AddInt64Value(value int64) *DataItem {
	if d.int64Value == nil {
		d.int64Value = new(int64)
	}
	*d.int64Value += value
	return d
}
func (d *DataItem) AddFloat64Value(value float64) *DataItem {
	if d.float64Value == nil {
		d.float64Value = new(float64)
	}

	*d.float64Value += value
	return d
}
func (d *DataItem) AddIntValue(value int) *DataItem {
	if d.intValue == nil {
		d.intValue = new(int)
	}
	*d.intValue += value
	return d
}

func (d *DataItem) Bytes() *DataItem {
	if d.int64Value != nil {
		d.Value = humanize.Bytes(uint64(*d.int64Value))
		return d
	}

	if d.float64Value != nil {
		d.Value = humanize.Bytes(uint64(*d.float64Value))
		return d
	}
	if d.intValue != nil {
		d.Value = humanize.Bytes(uint64(*d.intValue))
		return d
	}
	return d
}

func (d *DataItem) Number() *DataItem {
	if d.int64Value != nil {
		d.Value = humanize.Comma(*d.int64Value)
		return d
	}

	if d.float64Value != nil {
		d.Value = humanize.CommafWithDigits(*d.float64Value, 2)
		return d
	}

	if d.intValue != nil {
		d.Value = humanize.Comma(int64(*d.intValue))
		return d
	}
	return d
}

func (d *DataItem) Percent() *DataItem {

	if d.int64Value != nil {
		d.Value = humanize.Comma(*d.int64Value)

	}

	if d.float64Value != nil {
		d.Value = humanize.CommafWithDigits(*d.float64Value, 2)

	}

	if d.intValue != nil {
		d.Value = humanize.Comma(int64(*d.intValue))

	}
	d.Value = fmt.Sprintf("%s %", d.Value)
	return d
}

type DataItemsMap struct {
	keys []string
	m    map[string]*DataItem
}

func NewDataItemsMap(keys []string) *DataItemsMap {
	d := &DataItemsMap{
		keys: keys,
		m:    map[string]*DataItem{},
	}
	for _, key := range keys {
		d.m[key] = NewDataItem(key)
	}
	return d
}

func (d *DataItemsMap) Item(key string) *DataItem {
	return d.m[key]
}
func (d *DataItemsMap) List() []*DataItem {
	var list []*DataItem
	for _, key := range d.keys {
		list = append(list, d.m[key])
	}
	return list
}
