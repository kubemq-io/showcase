package api

import (
	"context"
	"fmt"
	"github.com/kubemq-io/showcase/apps/collector/pkg/types/api"
	"github.com/kubemq-io/showcase/apps/collector/pkg/types/base"
	kubemq2 "github.com/kubemq-io/showcase/apps/collector/pkg/types/kubemq"
	"github.com/kubemq-io/showcase/apps/collector/services/collector"
	"github.com/kubemq-io/showcase/apps/collector/services/kubemq"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"strconv"
	"time"
)

type Server struct {
	echoWebServer *echo.Echo
	collector     *collector.Collector
	kubemqService *kubemq.Service
}

func Start(ctx context.Context, collector *collector.Collector, kubemqService *kubemq.Service, port int) (*Server, error) {
	s := &Server{
		echoWebServer: echo.New(),
		collector:     collector,
		kubemqService: kubemqService,
	}
	s.echoWebServer.Use(middleware.Recover())
	s.echoWebServer.Use(middleware.CORS())
	s.echoWebServer.HideBanner = true
	s.echoWebServer.Static("/dashboard", "./dist")
	s.echoWebServer.GET("/health", func(c echo.Context) error {
		return c.String(200, "ok")
	})
	s.echoWebServer.GET("/ready", func(c echo.Context) error {

		return c.String(200, "ready")

	})
	s.echoWebServer.POST("/report", func(c echo.Context) error {
		m := &base.Metric{}
		err := c.Bind(m)
		if err != nil {
			return err
		}
		s.collector.Aggregate(m)
		return c.String(200, "")
	})
	s.echoWebServer.POST("/clear", func(c echo.Context) error {
		s.collector.ClearAll()
		return c.String(200, "")

	})
	s.echoWebServer.GET("/top", func(c echo.Context) error {
		return c.JSONPretty(200, s.collector.Top(c.QueryParam("group")), "\t")

	})
	s.echoWebServer.GET("/senders", func(c echo.Context) error {
		data := s.collector.Top("senders")
		return c.JSONPretty(200, api.GetSenders(data), "\t")

	})
	s.echoWebServer.GET("/receivers", func(c echo.Context) error {
		data := s.collector.Top("receivers")
		return c.JSONPretty(200, api.GetReceivers(data), "\t")
	})
	s.echoWebServer.GET("/api/status", func(c echo.Context) error {
		host := c.QueryParam("host")
		source := c.QueryParam("source")
		if source == "" {
			source = "queues"
		}
		res := kubemq2.NewResponse(c)
		res.SetResponseBody(kubemq.NewFakeStatus(host, source))
		return res.Send()
	})
	s.echoWebServer.GET("/kubemq", func(c echo.Context) error {
		return c.JSONPretty(200, api.GetKubeMQ(s.kubemqService.Get()), "\t")
	})
	s.echoWebServer.GET("/bucket", func(c echo.Context) error {
		count, _ := strconv.Atoi(c.QueryParam("count"))

		return c.JSONPretty(200, s.collector.Bucket(c.Param("source"), count), "\t")
	})
	errCh := make(chan error, 1)
	go func() {
		errCh <- s.echoWebServer.Start(fmt.Sprintf("0.0.0.0:%d", port))
	}()

	select {
	case err := <-errCh:
		if err != nil {
			return nil, err
		}
		return s, nil
	case <-time.After(1 * time.Second):
		return s, nil
	case <-ctx.Done():
		return nil, fmt.Errorf("error strarting api server, %w", ctx.Err())
	}
}

func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return s.echoWebServer.Shutdown(ctx)
}
