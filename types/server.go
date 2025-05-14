package types

import (
	"net/http/httputil"
	"sync"
)

type Server struct {
	sync.RWMutex
	Addr                string
	IsAlive             bool
	LastTimeOutResponse int
	Wieght              int
	Proxy               *httputil.ReverseProxy
}
