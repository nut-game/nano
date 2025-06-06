package main

import (
	"context"
	"flag"
	"fmt"
	"strconv"

	"strings"

	"github.com/nut-game/nano"
	"github.com/nut-game/nano/acceptor"
	"github.com/nut-game/nano/cluster"
	"github.com/nut-game/nano/component"
	"github.com/nut-game/nano/config"
	"github.com/nut-game/nano/constants"
	"github.com/nut-game/nano/examples/demo/cluster_grpc/services"
	"github.com/nut-game/nano/groups"
	"github.com/nut-game/nano/modules"
	"github.com/nut-game/nano/route"
)

var app nano.Nano

func configureBackend() {
	room := services.NewRoom(app)
	app.Register(room,
		component.WithName("room"),
		component.WithNameFunc(strings.ToLower),
	)

	app.RegisterRemote(room,
		component.WithName("room"),
		component.WithNameFunc(strings.ToLower),
	)
}

func configureFrontend(port int) {
	app.Register(services.NewConnector(app),
		component.WithName("connector"),
		component.WithNameFunc(strings.ToLower),
	)
	app.RegisterRemote(&services.ConnectorRemote{},
		component.WithName("connectorremote"),
		component.WithNameFunc(strings.ToLower),
	)

	err := app.AddRoute("room", func(
		ctx context.Context,
		route *route.Route,
		payload []byte,
		servers map[string]*cluster.Server,
	) (*cluster.Server, error) {
		// will return the first server
		for k := range servers {
			return servers[k], nil
		}
		return nil, nil
	})

	if err != nil {
		fmt.Printf("error adding route %s\n", err.Error())
	}

	err = app.SetDictionary(map[string]uint16{
		"connector.getsessiondata": 1,
		"connector.setsessiondata": 2,
		"room.room.getsessiondata": 3,
		"onMessage":                4,
		"onMembers":                5,
	})

	if err != nil {
		fmt.Printf("error setting route dictionary %s\n", err.Error())
	}
}

func main() {
	port := flag.Int("port", 3250, "the port to listen")
	svType := flag.String("type", "connector", "the server type")
	isFrontend := flag.Bool("frontend", true, "if server is frontend")
	rpcServerPort := flag.Int("rpcsvport", 3434, "the port that grpc server will listen")

	flag.Parse()

	meta := map[string]string{
		constants.GRPCHostKey: "127.0.0.1",
		constants.GRPCPortKey: strconv.Itoa(*rpcServerPort),
	}

	var bs *modules.ETCDBindingStorage
	app, bs = createApp(*port, *isFrontend, *svType, meta, *rpcServerPort)

	defer app.Shutdown()

	app.RegisterModule(bs, "bindingsStorage")
	if !*isFrontend {
		configureBackend()
	} else {
		configureFrontend(*port)
	}
	app.Start()
}

func createApp(port int, isFrontend bool, svType string, meta map[string]string, rpcServerPort int) (nano.Nano, *modules.ETCDBindingStorage) {
	builder := nano.NewDefaultBuilder(isFrontend, svType, nano.Cluster, meta, *config.NewDefaultNanoConfig())

	grpcServerConfig := builder.Config.Cluster.RPC.Server.Grpc
	grpcServerConfig.Port = rpcServerPort
	gs, err := cluster.NewGRPCServer(grpcServerConfig, builder.Server, builder.MetricsReporters)
	if err != nil {
		panic(err)
	}
	builder.RPCServer = gs
	builder.Groups = groups.NewMemoryGroupService(builder.Config.Groups.Memory)

	bs := modules.NewETCDBindingStorage(builder.Server, builder.SessionPool, builder.Config.Modules.BindingStorage.Etcd)

	gc, err := cluster.NewGRPCClient(
		builder.Config.Cluster.RPC.Client.Grpc,
		builder.Server,
		builder.MetricsReporters,
		bs,
		cluster.NewInfoRetriever(builder.Config.Cluster.Info),
	)
	if err != nil {
		panic(err)
	}
	builder.RPCClient = gc

	if isFrontend {
		tcp := acceptor.NewTCPAcceptor(fmt.Sprintf(":%d", port))
		builder.AddAcceptor(tcp)
	}

	return builder.Build(), bs
}
