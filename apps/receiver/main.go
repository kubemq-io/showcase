package main

import (
	"context"
	"fmt"
	"github.com/kubemq-io/kubemq-go"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func getClient(ctx context.Context, host string, port int, clientId string) (*kubemq.Client, error) {
	return kubemq.NewClient(ctx,
		kubemq.WithAddress(host, port),
		kubemq.WithClientId(clientId),
		kubemq.WithAutoReconnect(true),
		kubemq.WithReconnectInterval(time.Second),
		kubemq.WithTransportType(kubemq.TransportTypeGRPC))

}

func runQueueReceivers(ctx context.Context, cfg *Config, doneCh chan bool) {
	queueCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	appStats := CreateStats(cfg)
	var clients []StatsInterface
	for i := 1; i <= cfg.Receivers; i++ {
		log.Println(fmt.Sprintf("loading queue client %d", i))
		client := NewQueueClient(queueCtx, i+cfg.ChannelStartRange, cfg)
		clients = append(clients, client)
		time.Sleep(time.Duration(cfg.LoadInterval) * time.Millisecond)
		select {
		case <-doneCh:
			return
		default:

		}
	}
	for {
		select {
		case <-time.After(time.Duration(cfg.CollectEvery) * time.Second):
			appStats.CollectStats(clients).Print()
			go appStats.ReportStats()
		case <-doneCh:
			return
		case <-queueCtx.Done():
			return
		}
	}
}

func runStoreReceivers(ctx context.Context, cfg *Config, doneCh chan bool) {
	storeCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	appStats := CreateStats(cfg)
	var clients []StatsInterface
	for i := 1; i <= cfg.Receivers; i++ {
		log.Println(fmt.Sprintf("loading store client %d", i))
		client := NewStoreClient(storeCtx, i, cfg)
		clients = append(clients, client)
		time.Sleep(time.Duration(cfg.LoadInterval) * time.Millisecond)
		select {
		case <-doneCh:
			return
		default:

		}
	}
	for {
		select {
		case <-time.After(time.Duration(cfg.CollectEvery) * time.Second):
			appStats.CollectStats(clients).Print()
			go appStats.ReportStats()
		case <-doneCh:
			return
		case <-storeCtx.Done():
			return
		}
	}
}

func main() {
	var gracefulShutdown = make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, syscall.SIGTERM)
	signal.Notify(gracefulShutdown, syscall.SIGINT)
	signal.Notify(gracefulShutdown, syscall.SIGQUIT)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatalf(err.Error())
	}
	cfg.Print()
	doneCh := make(chan bool)
	switch cfg.Type {
	case "queue", "queues":
		go runQueueReceivers(ctx, cfg, doneCh)
	case "store", "st", "events_store":
		go runStoreReceivers(ctx, cfg, doneCh)
	default:
		fmt.Println("no valid type defined, aborting")
		return
	}
	if cfg.KillAfter > 0 {
		go func() {
			<-time.After(time.Duration(cfg.KillAfter) * time.Minute)
			fmt.Println("kill after timer expired")
			gracefulShutdown <- syscall.SIGTERM
		}()

	}

	select {
	case <-gracefulShutdown:
		doneCh <- true

	}
}
