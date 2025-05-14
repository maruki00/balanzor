package main

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"
)

func ConfigParse(pathCfg string) map[string]any {
	return map[string]any{}
}

func reverseRequest(u string) error {
	pu, err := url.Parse(u)
	if err != nil {
		return err
	}
	reverseProxy := httputil.NewSingleHostReverseProxy(pu)
	handler := func(proxy *httputil.ReverseProxy) func(rw http.ResponseWriter, r *http.Request) {
		return func(rw http.ResponseWriter, r *http.Request) {
			rw.Write([]byte("Done"))
			proxy.ServeHTTP(rw, r)
		}
	}
	http.HandleFunc("/lb", handler(reverseProxy))
	return http.ListenAndServe(":8082", nil)
}

func checkServerHealth(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	t := time.NewTicker(time.Second * 20)
	for {
		select {
		case <-t.C:
			for _, server := range servers {
				_ = server.CheckServerAlive(1)
				fmt.Printf("%#v\n", server)
			}

		case <-ctx.Done():
			return
		}
	}
}

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	slog.Info("loading ...")
	srvs := []string{
		"localhost:9090",
		"localhost:9091",
		"localhost:9092",
		"localhost:9093",
		"localhost:9094",
		"localhost:9095",
	}

	for _, srv := range srvs {
		servers = append(servers, Server{
			Addr:                srv,
			isAlive:             checkServerAlive(srv, 1),
			LastTimeOutResponse: math.MaxInt,
			Wieght:              0,
			Proxy:               nil,
		})
	}

	for _, srv := range servers {
		fmt.Printf("%#v\n", srv)
	}

	slog.Info("started")

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go checkServerHealth(ctx, wg)
	wg.Wait()

}
