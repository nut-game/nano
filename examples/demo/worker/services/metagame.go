package services

import (
	"context"

	"github.com/nut-game/nano/component"
	"github.com/nut-game/nano/examples/demo/worker/protos"
	"github.com/nut-game/nano/logger"
)

// Metagame server
type Metagame struct {
	component.Base
}

// LogRemote logs argument when called
func (m *Metagame) LogRemote(ctx context.Context, arg *protos.Arg) (*protos.Response, error) {
	logger.Log.Infof("argument %+v\n", arg)
	return &protos.Response{Code: 200, Msg: "ok"}, nil
}
