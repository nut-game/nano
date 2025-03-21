package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"strings"

	"github.com/nut-game/nano"
	"github.com/nut-game/nano/acceptor"
	"github.com/nut-game/nano/cluster"
	"github.com/nut-game/nano/component"
	"github.com/nut-game/nano/config"
	"github.com/nut-game/nano/examples/demo/cluster/services"
	"github.com/nut-game/nano/groups"
	"github.com/nut-game/nano/route"
	"github.com/nut-game/nano/tracing/jaeger"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	jaegercfg "github.com/uber/jaeger-client-go/config"
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

func configureFrontend() {
	app.Register(services.NewConnector(app),
		component.WithName("connector"),
		component.WithNameFunc(strings.ToLower),
	)

	app.RegisterRemote(services.NewConnectorRemote(app),
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

func configureJaeger(config *viper.Viper, logger logrus.FieldLogger) {
	cfg, err := jaegercfg.FromEnv()
	if cfg.ServiceName == "" {
		logger.Error("Could not init jaeger tracer without ServiceName, either set environment JAEGER_SERVICE_NAME or cfg.ServiceName = \"my-api\"")
		return
	}
	if err != nil {
		logger.Error("Could not parse Jaeger env vars: %s", err.Error())
		return
	}
	options := jaeger.Options{
		Disabled:    cfg.Disabled,
		Probability: cfg.Sampler.Param,
		ServiceName: cfg.ServiceName,
	}
	jaeger.Configure(options)
}

func main() {
	port := flag.Int("port", 3250, "the port to listen")
	svType := flag.String("type", "connector", "the server type")
	isFrontend := flag.Bool("frontend", true, "if server is frontend")

	flag.Parse()

	if os.Getenv("JAEGER_SERVICE_NAME") != "" {
		configureJaeger(viper.GetViper(), logrus.New())
	}

	builder := nano.NewDefaultBuilder(*isFrontend, *svType, nano.Cluster, map[string]string{}, *config.NewDefaultNanoConfig())
	if *isFrontend {
		tcp := acceptor.NewTCPAcceptor(fmt.Sprintf(":%d", *port))
		builder.AddAcceptor(tcp)
	}
	builder.Groups = groups.NewMemoryGroupService(builder.Config.Groups.Memory)
	app = builder.Build()

	defer app.Shutdown()

	if !*isFrontend {
		configureBackend()
	} else {
		configureFrontend()
	}

	app.Start()
}
