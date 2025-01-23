package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/nut-game/nano"
	"github.com/nut-game/nano/acceptor"
	"github.com/nut-game/nano/component"
	"github.com/nut-game/nano/config"
	"github.com/nut-game/nano/examples/demo/custom_metrics/services"
	"github.com/spf13/viper"
)

var app nano.Nano

func main() {
	port := flag.Int("port", 3250, "the port to listen")
	svType := "room"
	isFrontend := true

	flag.Parse()

	cfg := viper.New()
	cfg.AddConfigPath(".")
	cfg.SetConfigName("config")
	err := cfg.ReadInConfig()
	if err != nil {
		panic(err)
	}

	tcp := acceptor.NewTCPAcceptor(fmt.Sprintf(":%d", *port))

	conf := config.NewConfig(cfg)
	builder := nano.NewBuilderWithConfigs(isFrontend, svType, nano.Cluster, map[string]string{}, conf)
	builder.AddAcceptor(tcp)
	app = builder.Build()

	defer app.Shutdown()

	app.Register(services.NewRoom(app),
		component.WithName("room"),
		component.WithNameFunc(strings.ToLower),
	)

	app.Start()
}
