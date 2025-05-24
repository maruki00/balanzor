package types

import (
	"net"
	"net/http/httputil"
	"sync"
	"time"
)

type Option func(*Server)
type Server struct {
	sync.RWMutex
	Addr                string
	IsAlive             bool
	LastTimeOutResponse int
	Wieght              int
	Proxy               *httputil.ReverseProxy
}

func NewServer(options ...Option) *Server {
	s := Server{
		Addr:                "127.0.0.1",
		IsAlive:             false,
		LastTimeOutResponse: 0,
		Wieght:              0,
		Proxy:               nil,
	}

	for _, opt := range options {
		opt(&s)
	}
	return &s
}

func (_this *Server) CheckServerAlive(timeOut int) bool {
	_this.Lock()
	defer _this.Unlock()
	ATTEMPTS := 3

	_this.IsAlive = false
	for ATTEMPTS > 0 {
		start := time.Now()
		con, err := net.DialTimeout("tcp", _this.Addr, time.Duration(time.Second*time.Duration(timeOut)))
		responseTime := time.Since(start).Milliseconds()
		if err == nil {
			con.Close()
			_this.IsAlive = true
			_this.LastTimeOutResponse = int(responseTime)
			break
		}

		ATTEMPTS--
	}
	return _this.IsAlive
}

func WithAddress(addr string) Option {
	return func(s *Server) {
		s.Addr = addr
	}
}
func WithIsAlive(isAlive bool) Option {
	return func(s *Server) {
		s.IsAlive = isAlive
	}
}

func WithLastTimeOutResponse(lastTimeoutResponse int) Option {
	return func(s *Server) {
		s.LastTimeOutResponse = lastTimeoutResponse
	}
}

func WithWieght(wieght int) Option {
	return func(s *Server) {
		s.Wieght = wieght
	}
}

func WithProxy(proxy *httputil.ReverseProxy) Option {
	return func(s *Server) {
		s.Proxy = proxy
	}
}
