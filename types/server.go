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

func (_this *Server) checkServerAlive(timeOut int) bool {
	ATTEMPTS := 3
	for {
		conn, err := net.DialTimeout("tcp", _this.Addr, time.Duration(time.Second*time.Duration(timeOut)))
		conn.Close()
		if err == nil {
			return true
		}
		if ATTEMPTS <= 0 {
			return false
		}
		ATTEMPTS--
	}
	return true
}
