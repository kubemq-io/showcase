package kubemq

import (
	"github.com/kubemq-io/showcase/apps/collector/pkg/types/kubemq"
	"math/rand"
	"time"
)

func randomInt(min int, max int) int {
	return rand.Intn(max-min) + min
}

func randomInt64(min int, max int) int64 {
	return int64(rand.Intn(max-min) + min)
}
func randomFloat(min int, max int, divider int) float64 {
	return float64(rand.Intn(max-min)+min) / float64(divider)
}
func NewFakeStatus(host, source string) *kubemq.Status {
	s := kubemq.NewStatus(host)
	sys := kubemq.NewSystem()
	sys.CPUUtilization = randomFloat(0, 100, 100)
	sys.GoRoutines = randomInt64(100, 1000)
	sys.MemoryUtilization = randomFloat(0, 100, 100)
	sys.ProcessMemory = randomFloat(10000, 100000, 1)
	sys.TotalCPUs = randomInt(1, 64)
	s.SetSystem(sys)
	s.Entities = kubemq.NewEntitiesGroup()
	s.Entities[source] = kubemq.NewGroup().
		SetTotal(randomInt(1, 1000)).
		SetOut(kubemq.NewBaseValues().
			SetMessages(randomInt64(10000, 1000000)).
			SetVolume(randomInt64(10000, 1000000)).
			SetErrors(randomInt64(10000, 1000000)).
			SetWaiting(randomInt64(10000, 1000000))).
		SetIn(kubemq.NewBaseValues().
			SetMessages(randomInt64(10000, 1000000)).
			SetVolume(randomInt64(10000, 1000000)).
			SetErrors(randomInt64(10000, 1000000)))

	return s
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
