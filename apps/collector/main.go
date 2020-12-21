package main

import (
	"context"
	"github.com/kubemq-io/showcase/apps/collector/pkg/logger"
	"github.com/kubemq-io/showcase/apps/collector/services/api"
	"github.com/kubemq-io/showcase/apps/collector/services/collector"
	"github.com/kubemq-io/showcase/apps/collector/services/console"
	"os"
	"os/signal"
	"syscall"
)

var log *logger.Logger

func run() error {
	var gracefulShutdown = make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, syscall.SIGTERM)
	signal.Notify(gracefulShutdown, syscall.SIGINT)
	signal.Notify(gracefulShutdown, syscall.SIGQUIT)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	col, err := collector.NewCollector(ctx)
	if err != nil {
		return err
	}

	apiServer, err := api.Start(ctx, col, 8085)
	if err != nil {
		return err
	}
	_ = console.NewConsole(ctx, col)
	<-gracefulShutdown
	_ = apiServer.Stop()
	return nil

}

func main() {
	log = logger.NewLogger("main")
	log.Info("starting collector")
	if err := run(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}
