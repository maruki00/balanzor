package main

import "balazor/types"

type RoundRoubin struct {
	Backends      []*types.Server
	ServersLenght int
	CurrentNode   int
}

func (_this *RoundRoubin) GetNextNode() *types.Server {
	for i := _this.CurrentNode; i < _this.CurrentNode+_this.ServersLenght; i++ {
		if !_this.Backends[i%_this.ServersLenght].IsAlive {
			continue
		}
		_this.CurrentNode = (i % _this.ServersLenght) + 1
		return _this.Backends[i%_this.ServersLenght]
	}
	return nil
}
