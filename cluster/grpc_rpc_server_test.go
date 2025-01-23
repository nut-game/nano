package cluster

import (
	"fmt"
	"net"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/nut-game/nano/config"
	"github.com/nut-game/nano/helpers"
	"github.com/nut-game/nano/metrics"
	protosmocks "github.com/nut-game/nano/protos/mocks"
	"github.com/stretchr/testify/assert"
)

func TestNewGRPCServer(t *testing.T) {
	t.Parallel()
	sv := getServer()
	gs, err := NewGRPCServer(config.NewDefaultNanoConfig().Cluster.RPC.Server.Grpc, sv, []metrics.Reporter{})
	assert.NoError(t, err)
	assert.NotNil(t, gs)
}

func TestGRPCServerInit(t *testing.T) {
	t.Parallel()
	c := config.NewDefaultNanoConfig().Cluster.RPC.Server.Grpc
	c.Port = helpers.GetFreePort(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockNanoServer := protosmocks.NewMockNanoServer(ctrl)

	sv := getServer()
	gs, err := NewGRPCServer(c, sv, []metrics.Reporter{})
	gs.SetNanoServer(mockNanoServer)
	err = gs.Init()
	assert.NoError(t, err)
	assert.NotNil(t, gs)

	conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", c.Port))
	assert.NoError(t, err)
	assert.NotNil(t, conn)
	assert.NotNil(t, gs.grpcSv)
}
