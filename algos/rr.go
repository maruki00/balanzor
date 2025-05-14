package algos

import (
	"balazor/types"
	"context"
	"time"
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

func (_this *RoundRoubin) GetNextNode() *types.Server {
	for i := _this.CurrentNode; i < _this.CurrentNode+_this.ServersLenght; i++ {
		if !_this.Servers[i%_this.ServersLenght].IsAlive {
			continue
		}
		_this.CurrentNode = (i % _this.ServersLenght) + 1
		node := _this.Servers[i%_this.ServersLenght]
		node.Wieght++
		return node
	}
	return nil
}

func (_this *RoundRoubin) CheckServersHealth(ctx context.Context) {
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
