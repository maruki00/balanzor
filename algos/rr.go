package main

import (
	"balazor/types"
	"context"
	"fmt"
	"sync"
	"time"
)

type RoundRoubin struct {
	Servers       []*types.Server
	ServersLenght int
	CurrentNode   int
}

func (_this *RoundRoubin) GetNextNode() *types.Server {
	for i := _this.CurrentNode; i < _this.CurrentNode+_this.ServersLenght; i++ {
		if !_this.Servers[i%_this.ServersLenght].IsAlive {
			continue
		}
		_this.CurrentNode = (i % _this.ServersLenght) + 1
		return _this.Backends[i%_this.ServersLenght]
	}
	return nil
}

func (_this *RoundRoubin) checkServerHealth(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	t := time.NewTicker(time.Second * 20)
	for {
		select {
		case <-t.C:
			for _, srv := range _this.Servers {
				_ = srv.CheckServerAlive(1)
				fmt.Printf("%#v\n", srv)
			}

		case <-ctx.Done():
			return
		}
	}
}
