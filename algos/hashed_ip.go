package algos

import (
	"balazor/types"
	"context"
	"hash/fnv"
	"time"
)

type HashedIP struct {
	Servers       []*types.Server
	ServersLenght int
}

func (_this *HashedIP) AppendServer(server *types.Server) {
	_this.Servers = append(_this.Servers, server)
	_this.ServersLenght++
}

func (_this *HashedIP) SetServers(servers []*types.Server) {
	copy(_this.Servers, servers)
	_this.ServersLenght = len(servers)
}

func (_this *HashedIP) GetServers() []*types.Server {
	return _this.Servers
}

func (_this *HashedIP) GetServer(index int) *types.Server {
	return _this.Servers[index]
}

func (_this *HashedIP) GetCurrentNode(ctx types.BalanzerCtx) *types.Server {
	ip := ctx.IP
	currNode := _this.Servers[_this.hashIP(ip)]
	return currNode
}

func (_this *HashedIP) CheckServersHealth(ctx context.Context) {
	t := time.NewTicker(time.Second * 1)
	for {
		select {
		case <-t.C:
			for _, srv := range _this.GetServers() {
				_ = srv.CheckServerAlive(1)
			}
		case <-ctx.Done():
			return
		}
	}
}

// utils

func (_this *HashedIP) hashIP(ip string) int {
	hash := fnv.New32()
	hash.Write([]byte(ip))
	return int(hash.Sum32()) % _this.ServersLenght
}
