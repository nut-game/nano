package services

import (
	"context"

	pitaya "github.com/nut-game/nano"
	"github.com/nut-game/nano/component"
	"github.com/nut-game/nano/examples/demo/worker/protos"
	"google.golang.org/protobuf/proto"
)

// Worker server
type Worker struct {
	component.Base
}

// Configure starts workers and register rpc job
func (w *Worker) Configure(app pitaya.Pitaya) {
	app.StartWorker()
	app.RegisterRPCJob(&RPCJob{app: app})
}

// RPCJob implements worker.RPCJob
type RPCJob struct {
	app pitaya.Pitaya
}

// ServerDiscovery returns a serverID="", meaning any server
// is ok
func (r *RPCJob) ServerDiscovery(
	route string,
	rpcMetadata map[string]interface{},
) (serverID string, err error) {
	return "", nil
}

// RPC calls pitaya's rpc
func (r *RPCJob) RPC(
	ctx context.Context,
	serverID, routeStr string,
	reply, arg proto.Message,
) error {
	return r.app.RPCTo(ctx, serverID, routeStr, reply, arg)
}

// GetArgReply returns reply and arg of LogRemote,
// since we have no other methods in this example
func (r *RPCJob) GetArgReply(
	route string,
) (arg, reply proto.Message, err error) {
	return &protos.Arg{}, &protos.Response{}, nil
}
