package main

import "balazor/types"

type RoundRoubin struct {
	Backends []types.Server
	ServersLenght int
	CurrentNode int
}


func (_this *RoundRoubin) GetNext() *types.Server {
	for i:= _this.CurrentNode; i< _this.CurrentNode+_this.ServersLenght; i++ {
		if !_this.Backends[i%_this.ServersLenght].IsLive
	} 
}
