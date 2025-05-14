package algos

import (
	"balazor/types"
	"context"
)

type Algo interface {
	GetCurrentNode() *types.Server
	CheckServersHealth(context.Context)
	AppendServer(*types.Server)
	SetServers([]*types.Server)
	GetServers() []*types.Server
	GetServer(int) *types.Server
}
