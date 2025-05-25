package algos

import (
	"balazor/types"
	"context"
	"time"
)

type WeightedRoundRoubin struct {
	Servers       []*types.Server
	ServersLenght int
	CurrentNode   int
	MaxWeightt    int
}

func (_this *WeightedRoundRoubin) AppendServer(server *types.Server) {
	_this.Servers = append(_this.Servers, server)
	_this.ServersLenght++
}

func (_this *WeightedRoundRoubin) SetServers(servers []*types.Server) {
	copy(_this.Servers, servers)
	_this.ServersLenght = len(servers)
}

func (_this *WeightedRoundRoubin) GetServers() []*types.Server {
	return _this.Servers
}

func (_this *WeightedRoundRoubin) GetServer(index int) *types.Server {
	return _this.Servers[index]
}

func (_this *WeightedRoundRoubin) GetCurrentNode() *types.Server {
	var best *types.Server
	total := 0

	for _, server := range _this.Servers {
		total += server.Weight
		server.Weight += server.Weight
		if best == nil || server.Weight > best.Weight {
			best = server
		}
	}

	if best == nil {
		return nil
	}

	best.Weight -= total
	return best
}

func (_this *WeightedRoundRoubin) CheckServersHealth(ctx context.Context) {
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
