package algos

import (
	"balazor/types"
	"context"
	"sync"
)

type Algo interface {
	GetNextNode() *Types.Server
	CheckServersHealth(context.Context, *sync.WaitGroup)
	AppendServer(types.Server)
	SetServers([]*types.Server)
	GetServers() []*types.Server
	GetServer(int) *types.Server
}
