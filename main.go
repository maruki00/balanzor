package main

import (
	"balazor/algos"
	"balazor/types"
	"context"
	"log/slog"
	"math"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

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

func main() {
	algo := "round-roubin"
	// timeOut := 1
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
	var lb algos.Algo
	switch algo {
	case "round-roubin":
		lb = &algos.RoundRoubin{}
		break
	default:
		panic("algo not supported")
	}
	for _, srv := range srvs {
		srv := &types.Server{
			Addr:                srv,
			isAlive:             false,
			LastTimeOutResponse: math.MaxInt,
			Wieght:              0,
			Proxy:               nil,
		}
		lb.AppendServer(srv)
	}

	slog.Info("started")
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go lb.CheckServersHealth(ctx, wg)
	wg.Wait()
}
