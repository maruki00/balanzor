package algos

import (
	"context"

	"github.com/maruki00/balazor/types"
)

type Algo interface {
	GetCurrentNode(types.BalanzerCtx) *types.Server
	CheckServersHealth(context.Context)
	AppendServer(*types.Server)
	SetServers([]*types.Server)
	GetServers() []*types.Server
	GetServer(int) *types.Server
}
