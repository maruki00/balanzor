package main

import (
	"balazor/algos"
	"balazor/types"
	"context"
	"errors"
	"log/slog"
	"math"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func reverseRequest(lb algos.Algo, rw http.ResponseWriter, r *http.Request) error {
	curNode := lb.GetNextNode()
	if curNode == nil {
		return errors.New("no Server is alive")
	}
	pu, err := url.Parse(curNode.Addr)
	if err != nil {
		panic(err)
	}

	reverseProxy := httputil.NewSingleHostReverseProxy(pu)
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

	slog.Info("started")

	go lb.CheckServersHealth(ctx)
	http.HandleFunc("/lb", func(rw http.ResponseWriter, r *http.Request) {
		err := reverseRequest(lb, rw, r)
		if err != nil {
			rw.Write([]byte(err.Error()))
		}
	})

	slog.Info("Start Server 0.0.0.0:8082")
	http.ListenAndServe(":8082", nil)

}
