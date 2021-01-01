package kubemq

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/kubemq-io/showcase/apps/collector/pkg/types/kubemq"
	"log"
	"strings"
	"sync"
	"time"
)

const reportInterval = 5 * time.Second
const statusEndPoint = "/api/status"

type Service struct {
	sync.Mutex
	collectHosts []string
	statusMap    map[string]*kubemq.Status
	restyClient  *resty.Client
}

func NewService(hosts string) (*Service, error) {
	s := &Service{
		collectHosts: strings.Split(hosts, ","),
		statusMap:    map[string]*kubemq.Status{},
		restyClient:  resty.New(),
	}

	return s, nil
}
func (s *Service) Start(ctx context.Context) {
	go s.run(ctx)
}
func (s *Service) getStatus(ctx context.Context, host string) {
	response := &kubemq.Response{}
	url := fmt.Sprintf("http://%s%s?host=%s&source=queues", host, statusEndPoint, host)
	_, err := s.restyClient.R().SetContext(ctx).SetResult(response).Get(url)
	if err != nil {
		log.Println(fmt.Sprintf("error get kubemq status from: %s", url))
		return
	}
	if len(response.Data) == 0 {
		log.Println(fmt.Sprintf("error get kubemq status from: %s, status empty", url))
		return
	}
	status := &kubemq.Status{}
	err = json.Unmarshal(response.Data, status)
	if err != nil {
		log.Println(fmt.Sprintf("error unmarshal kubemq status from: %s error: %s", url, err.Error()))
		return
	}
	s.Lock()
	defer s.Unlock()
	s.statusMap[host] = status
}
func (s *Service) run(ctx context.Context) {
	for {
		select {
		case <-time.After(reportInterval):
			for _, host := range s.collectHosts {
				s.getStatus(ctx, host)
			}
		case <-ctx.Done():
			return
		}
	}
}

func (s *Service) Get() map[string]*kubemq.Status {
	s.Lock()
	defer s.Unlock()
	resp := s.statusMap
	return resp
}
func (s *Service) Clear() {
	s.Lock()
	defer s.Unlock()
	s.statusMap = map[string]*kubemq.Status{}
}
