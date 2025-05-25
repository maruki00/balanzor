package main

import (
	"balazor/algos"
	"balazor/types"
	"context"
	"fmt"
	"log"
	"log/slog"
	"math"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

var (
	lb algos.Algo
)

func balancer(rw http.ResponseWriter, r *http.Request) {
	go reverseRequest(rw, r)
}

func reverseRequest(rw http.ResponseWriter, r *http.Request) {
	curNode := lb.GetCurrentNode()
	if curNode == nil {
		rw.Write([]byte("server not available."))
		return
	}
	fmt.Println(curNode.Addr, " --> ", curNode.Wieght)
	curNode.Proxy.ServeHTTP(rw, r)
	curNode.Wieght--
}

func main() {

	cfg := types.NewConfig("config.yaml")
	fmt.Println(cfg)
	return

	algo := "round-roubin"
	// timeOut := 1
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	slog.Info("loading ...")
	srvs := []string{
		"http://127.0.0.1:9090",
		"http://127.0.0.1:9091",
		"http://127.0.0.1:9092",
		"http://127.0.0.1:9093",
		"http://127.0.0.1:9094",
		"http://127.0.0.1:9095",
	}
	switch algo {
	case "round-roubin":
		lb = &algos.RoundRoubin{}
	default:
		panic("algo not supported")
	}
	for _, srv := range srvs {

		srvUri, err := url.Parse(srv)
		if err != nil {
			panic("error : " + err.Error())
		}

		srv := &types.Server{
			Addr:                srvUri.Host,
			IsAlive:             false,
			LastTimeOutResponse: math.MaxInt,
			Wieght:              0,
			Proxy:               nil,
		}

		slog.Info("uri", srvUri.String(), srvUri.Host)

		proxy := httputil.NewSingleHostReverseProxy(srvUri)
		proxy.ErrorHandler = func(writer http.ResponseWriter, request *http.Request, e error) {
			log.Printf("[%s] %s\n", srvUri.Host, e.Error())

			<-time.After(10 * time.Millisecond)
			proxy.ServeHTTP(writer, request.WithContext(ctx))

			balancer(writer, request.WithContext(ctx))
		}
		srv.Proxy = proxy
		lb.AppendServer(srv)
	}

	slog.Info("started")

	go lb.CheckServersHealth(ctx)
	http.HandleFunc("/lb", func(writer http.ResponseWriter, request *http.Request) {
		reverseRequest(writer, request)
	})

	slog.Info("Start Server 0.0.0.0:8082")
	http.ListenAndServe(":8082", nil)

}
