package api

import (
	"context"
	"fmt"
	"github.com/kubemq-io/showcase/apps/collector/pkg/types"
	"github.com/kubemq-io/showcase/apps/collector/services/collector"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"strconv"
	"time"
)

type Server struct {
	echoWebServer *echo.Echo
	collector     *collector.Collector
}

func Start(ctx context.Context, collector *collector.Collector, port int) (*Server, error) {
	s := &Server{
		echoWebServer: echo.New(),
		collector:     collector,
	}
	s.echoWebServer.Use(middleware.Recover())
	s.echoWebServer.Use(middleware.CORS())
	s.echoWebServer.HideBanner = true

	s.echoWebServer.GET("/health", func(c echo.Context) error {
		return c.String(200, "ok")
	})
	s.echoWebServer.GET("/ready", func(c echo.Context) error {

		return c.String(200, "ready")

	})
	s.echoWebServer.POST("/report", func(c echo.Context) error {
		m := &types.Metric{}
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
		return c.JSONPretty(200, s.collector.Top(), "\t")

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
