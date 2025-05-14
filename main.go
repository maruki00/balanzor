package main

import (
	"balazor/algos"
	"balazor/types"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"net/http"
	"net/http/httputil"
	"sync"
)

func reverseRequest(lb algos.Algo, rw http.ResponseWriter, r *http.Request) error {
	curNode := lb.GetNextNode()
	if curNode == nil {
		return errors.New("No Server is alive")
	}
	pu. err := 
	reverseProxy := httputil.NewSingleHostReverseProxy()
	reverseProxy.ServeHTTP(rw, r)	
	return nil
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
			IsAlive:             false,
			LastTimeOutResponse: math.MaxInt,
			Wieght:              0,
			Proxy:               nil,
		}
		lb.AppendServer(srv)
	}
	//reverseRequest(lb)

	http.HandleFunc("/lb", func(rw http.ResponseWriter, r *http.Request) {
		fmt.Println("hello world")
		reverseRequest(lb, rw, r)
	})

	slog.Info("Start Server 0.0.0.0:8082")
	http.ListenAndServe(":8082", nil)

	slog.Info("started")
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go lb.CheckServersHealth(ctx, wg)
	wg.Wait()
}
