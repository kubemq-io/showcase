package main

import (
	"sync"
	"time"
)

type StorageRec struct {
	Channel   string
	MessageID string
	Timestamp time.Time
}

type Storage struct {
	channel     string
	addMutex    sync.Mutex
	deleteMutex sync.Mutex
	addMap      map[string]*StorageRec
	deleteMap   map[string]string
}

func InitStorage(channel string) *Storage {
	return &Storage{
		channel:     channel,
		addMutex:    sync.Mutex{},
		deleteMutex: sync.Mutex{},
		addMap:      map[string]*StorageRec{},
		deleteMap:   map[string]string{},
	}
}

func (s *Storage) Add(msgId string) {
	s.addMutex.Lock()
	defer s.addMutex.Unlock()
	s.addMap[msgId] = &StorageRec{
		Channel:   s.channel,
		MessageID: msgId,
		Timestamp: time.Now(),
	}
}

func (s *Storage) Delete(msgId string) {
	s.deleteMutex.Lock()
	defer s.deleteMutex.Unlock()
	s.deleteMap[msgId] = msgId

}

func (s *Storage) cleanQueues() {
	s.addMutex.Lock()
	defer s.addMutex.Unlock()
	s.deleteMutex.Lock()
	defer s.deleteMutex.Unlock()
	deleteList := []string{}
	for msgId, _ := range s.deleteMap {
		_, ok := s.addMap[msgId]
		if ok {
			s.addMap[msgId] = nil
			delete(s.addMap, msgId)
			deleteList = append(deleteList, msgId)
		}
	}
	for _, msgID := range deleteList {
		delete(s.deleteMap, msgID)
	}
	deleteList = nil
}

func (s *Storage) StillInQueue() int {
	s.cleanQueues()
	s.addMutex.Lock()
	defer s.addMutex.Unlock()
	return len(s.addMap)

}
func (s *Storage) StillInDeleteQueue() int {
	s.deleteMutex.Lock()
	defer s.deleteMutex.Unlock()
	return len(s.deleteMap)

}

//func (s *Storage) CheckSeq() int {
//	s.deleteMutex.Lock()
//	defer s.deleteMutex.Unlock()
//	cnt := 0
//	if len(s.seqList) > 1 {
//		for i := 1; i < len(s.seqList); i++ {
//			item := s.seqList[i]
//			itemBefore := s.seqList[i-1]
//			if item-itemBefore != 1 {
//				//		log.Println(fmt.Sprintf("bad seq on channel: %s, item: %d, item before: %d", s.channel, item, itemBefore))
//				cnt++
//			}
//		}
//	}
//	return cnt
//}
