package main

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var done = make(chan int)

func main() {
	group, ctx := errgroup.WithContext(context.Background())

	group.Go(func() error {
		mux := http.NewServeMux()
		mux.HandleFunc("/close", func(writer http.ResponseWriter, request *http.Request) {
			done <- 1
		})

		s := NewServer(":8081", mux)
		go func() {
			err := s.Start()
			if err != nil {
				fmt.Println(err)
			}
		}()

		select {
		case <-done:
			return s.Stop()
		case <-ctx.Done():
			return errors.New("信号关闭")
		}
	})

	group.Go(func() error {
		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)

		select {
		case <-quit:
			return errors.New("信号关闭")
		case <-ctx.Done():
			return errors.New("httpServer 关闭")
		}
	})

	// 捕获err
	fmt.Println("开始捕获error")
	err := group.Wait()
	fmt.Println("关闭方式为:", err)
}

// httpServer 结构体
type httpServer struct {
	server  *http.Server
	handler http.Handler
	cxt     context.Context
}

// NewServer 定义一个新的http server
func NewServer(address string, mux http.Handler) *httpServer {
	// 定义cxt
	h := &httpServer{cxt: context.Background()}
	// 定义server
	h.server = &http.Server{
		Addr:         address,
		WriteTimeout: time.Second * 3,
		Handler:      mux,
	}
	return h
}

// Start 服务开启
func (h *httpServer) Start() error {
	fmt.Println("httpServer 开始")
	return h.server.ListenAndServe()
}

// Stop 服务关闭
func (h *httpServer) Stop() error {
	_ = h.server.Shutdown(h.cxt)
	return fmt.Errorf("httpServer 关闭")
}