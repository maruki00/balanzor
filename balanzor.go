package balanzor

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"math"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/maruki00/balanzor/algos"
	"github.com/maruki00/balanzor/types"
)

var (
	lb algos.Algo
)

type Balanazor struct {
	backends         []string
	algorithm        string
	addr             string
	endpoint         string
	healthCheckDelay int
}

func NewBalanzor(
	backends []string,
	algorithm string,
	addr string,
	endpoint string,
	delay int,
) *Balanazor {
	return &Balanazor{
		backends:         backends,
		algorithm:        algorithm,
		addr:             addr,
		endpoint:         endpoint,
		healthCheckDelay: delay,
	}
}

func (_this *Balanazor) balancer(rw http.ResponseWriter, r *http.Request) {
	go _this.reverseRequest(rw, r)
}

func (_this *Balanazor) reverseRequest(rw http.ResponseWriter, r *http.Request) {
	curNode := lb.GetCurrentNode(types.BalanzerCtx{
		Ctx: context.TODO(),
		IP:  r.Host,
	})
	if curNode == nil {
		rw.Write([]byte("server not available."))
		return
	}
	fmt.Println(curNode.Addr, " --> ", curNode.Weight)
	curNode.Proxy.ServeHTTP(rw, r)
	curNode.Weight--
}

func (_this *Balanazor) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	slog.Info("loading ...", _this.algorithm, " algorithm")
	switch _this.algorithm {
	case "round-roubin":
		lb = &algos.RoundRoubin{}
	case "weighted-round-roubin":
		lb = &algos.WeightedRoundRoubin{}
	case "hashed-ip":
		lb = &algos.HashedIP{}
	default:
		return errors.New("algo not supported")
	}
	for _, srv := range _this.backends {
		srvUri, err := url.Parse(srv)
		if err != nil {
			return err
		}
		srv := &types.Server{
			Addr:                srvUri.Host,
			IsAlive:             false,
			LastTimeOutResponse: math.MaxInt,
			Weight:              0,
			Proxy:               nil,
		}
		slog.Info("uri", srvUri.String(), srvUri.Host)
		proxy := httputil.NewSingleHostReverseProxy(srvUri)
		proxy.ErrorHandler = func(writer http.ResponseWriter, request *http.Request, e error) {
			log.Printf("[%s] %s\n", srvUri.Host, e.Error())
			<-time.After(10 * time.Millisecond)
			proxy.ServeHTTP(writer, request.WithContext(ctx))
			_this.balancer(writer, request.WithContext(ctx))
		}
		srv.Proxy = proxy
		lb.AppendServer(srv)
	}
	slog.Info("started")
	go lb.CheckServersHealth(ctx)
	http.HandleFunc(_this.endpoint, func(writer http.ResponseWriter, request *http.Request) {
		_this.reverseRequest(writer, request)
	})
	slog.Info("Start Server ", "server", _this.addr)
	return http.ListenAndServe(_this.addr, nil)

}
