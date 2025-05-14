package types

import (
	"net"
	"net/http/httputil"
	"sync"
	"time"
)

type Server struct {
	sync.RWMutex
	Addr                string
	IsAlive             bool
	LastTimeOutResponse int
	Wieght              int
	Proxy               *httputil.ReverseProxy
}

func (_this *Server) CheckServerAlive(timeOut int) bool {
	ATTEMPTS := 3
	_this.IsAlive = false
	for ATTEMPTS > 0 {
		start := time.Now()
		_, err := net.DialTimeout("tcp", _this.Addr, time.Duration(time.Second*time.Duration(timeOut)))
		responseTime := time.Since(start).Milliseconds()
		if err == nil {
			_this.IsAlive = true
			_this.LastTimeOutResponse = int(responseTime)
			break
		}
		ATTEMPTS--
	}
	return _this.IsAlive
}
