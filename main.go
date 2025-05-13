package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"
)

var (
	servers []Server
)

type Server struct {
	Addr                string
	isAlive             bool
	LastTimeOutResponse int
	Wieght              int
	Proxy               *httputil.ReverseProxy
}

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

func chechServerAlive(u string, timeOut int) bool {
	_, err := net.DialTimeout("udp", u, time.Duration(time.Second*time.Duration(timeOut)))
	return err != nil
}

func checkServerHealth(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	t := time.NewTicker(time.Second * 2)
	for {
		select {
		case <-t.C:
			for i := range 10 {
				fmt.Println("check server: ", i)
			}
		case <-ctx.Done():
			break
		}
	}

}

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go checkServerHealth(ctx, wg)
	wg.Wait()

}
