rpackage algos

import (
	"balazor/types"
	"context"
	"time"
)

type WieghtedRoundRoubin struct {
	Servers       []*types.Server
	ServersLenght int
	CurrentNode   int
	MaxWieghtt    int
}

func (_this *WieghtedRoundRoubin) AppendServer(server *types.Server) {
	_this.Servers = append(_this.Servers, server)
	_this.ServersLenght++
}

func (_this *WieghtedRoundRoubin) SetServers(servers []*types.Server) {
	copy(_this.Servers, servers)
	_this.ServersLenght = len(servers)
}

func (_this *WieghtedRoundRoubin) GetServers() []*types.Server {
	return _this.Servers
}

func (_this *WieghtedRoundRoubin) GetServer(index int) *types.Server {
	return _this.Servers[index]
}

func (_this *WieghtedRoundRoubin) GetCurrentNode() *types.Server {
	for i := range _this.ServersLenght {
		currIndex := (i + _this.CurrentNode) % _this.ServersLenght
		if !_this.Servers[currIndex].IsAlive || !_this.Servers[currIndex].CheckServerAlive(1) {
			continue
		}
		if _this.Servers[currIndex].Wieght >= _this.MaxWieghtt {
			continue
		}

		_this.CurrentNode = currIndex + 1

		node := _this.Servers[currIndex]
		node.Wieght++
		return node
	}
	return nil
}

func (_this *WieghtedRoundRoubin) CheckServersHealth(ctx context.Context) {
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
