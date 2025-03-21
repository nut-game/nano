package networkentity

import (
	"context"
	"net"

	"github.com/nut-game/nano/protos"
)

// NetworkEntity represent low-level network instance
type NetworkEntity interface {
	Push(route string, v any) error
	ResponseMID(ctx context.Context, mid uint, v any, isError ...bool) error
	Close() error
	Kick(ctx context.Context) error
	RemoteAddr() net.Addr
	SendRequest(ctx context.Context, serverID, route string, v any) (*protos.Response, error)
}
