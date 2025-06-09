package algos

import (
	"context"
	"time"

	"github.com/maruki00/balanzor/types"
)

type RoundRoubin struct {
	Servers       []*types.Server
	ServersLenght int
	CurrentNode   int
}

func (_this *RoundRoubin) AppendServer(server *types.Server) {
	_this.Servers = append(_this.Servers, server)
	_this.ServersLenght++
}

func (_this *RoundRoubin) SetServers(servers []*types.Server) {
	copy(_this.Servers, servers)
	_this.ServersLenght = len(servers)
}

func (_this *RoundRoubin) GetServers() []*types.Server {
	return _this.Servers
}

func (_this *RoundRoubin) GetServer(index int) *types.Server {
	return _this.Servers[index]
}

func (_this *RoundRoubin) GetCurrentNode(ctx types.BalanzerCtx) *types.Server {
	for i := range _this.ServersLenght {
		currIndex := (i + _this.CurrentNode) % _this.ServersLenght
		if !_this.Servers[currIndex].IsAlive || !_this.Servers[currIndex].CheckServerAlive(1) {
			continue
		}
		_this.CurrentNode = currIndex + 1
		node := _this.Servers[currIndex]
		node.Weight++
		return node
	}
	return nil
}

func (_this *RoundRoubin) CheckServersHealth(ctx context.Context) {
	t := time.NewTicker(time.Second * 10)
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
