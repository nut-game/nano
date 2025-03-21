package services

import (
	"context"

	"github.com/nut-game/nano"
	"github.com/nut-game/nano/component"
	"github.com/nut-game/nano/examples/demo/worker/protos"
	"github.com/nut-game/nano/logger"
)

// Room server
type Room struct {
	component.Base
	app nano.Nano
}

// NewRoom ctor
func NewRoom(app nano.Nano) *Room {
	return &Room{app: app}
}

// CallLog makes ReliableRPC to metagame LogRemote
func (r *Room) CallLog(ctx context.Context, arg *protos.Arg) (*protos.Response, error) {
	route := "metagame.metagame.logremote"
	reply := &protos.Response{}
	jid, err := r.app.ReliableRPC(route, nil, reply, arg)
	if err != nil {
		logger.Infof("failed to enqueue rpc: %q", err)
		return nil, err
	}

	logger.Infof("enqueue rpc job: %d", jid)
	return &protos.Response{Code: 200, Msg: "ok"}, nil
}
