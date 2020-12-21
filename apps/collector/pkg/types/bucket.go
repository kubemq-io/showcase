package types

import "sync"

type Bucket struct {
	sync.Mutex
	Name  string
	Items []*Snapshot
}

func NewBucket(name string) *Bucket {
	return &Bucket{
		Mutex: sync.Mutex{},
		Name:  name,
		Items: []*Snapshot{},
	}
}

func (b *Bucket) Append(s *Snapshot) {
	b.Lock()
	defer b.Unlock()
	b.Items = append([]*Snapshot{s}, b.Items...)
}

func (b *Bucket) Top() *Snapshot {
	b.Lock()
	defer b.Unlock()
	if len(b.Items) > 0 {
		return b.Items[0]
	}
	return nil
}
func (b *Bucket) List(count int) []*Snapshot {
	b.Lock()
	defer b.Unlock()

	if count == 0 || count >= len(b.Items) {
		return b.Items
	} else {
		return b.Items[:count]
	}
}
