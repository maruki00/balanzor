package types

import "context"

type BalanzerCtx struct {
	context.Context
	IP string
}
