package web

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Server struct {
}

func StartServer(ctx context.Context, port int, path, staticFiles string) error {
	fs := http.FileServer(http.Dir(staticFiles))
	http.Handle(path, fs)
	errCh := make(chan error, 1)
	go func() {
		errCh <- http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	}()
	select {
	case err := <-errCh:
		if err != nil {
			return err
		}
		return nil
	case <-time.After(1 * time.Second):
		return nil
	case <-ctx.Done():
		return fmt.Errorf("error strarting web server, %w", ctx.Err())
	}
}
