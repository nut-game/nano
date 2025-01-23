package main

import (
	"flag"
	"fmt"

	"strings"

	"github.com/nut-game/nano"
	"github.com/nut-game/nano/acceptor"
	"github.com/nut-game/nano/component"
	"github.com/nut-game/nano/config"
	"github.com/nut-game/nano/examples/demo/worker/services"
	"github.com/spf13/viper"
)

var app nano.Nano

func configureWorker() {
	worker := services.Worker{}
	worker.Configure(app)
}

func main() {
	port := flag.Int("port", 3250, "the port to listen")
	svType := flag.String("type", "metagame", "the server type")
	isFrontend := flag.Bool("frontend", true, "if server is frontend")

	flag.Parse()

	conf := viper.New()
	conf.SetDefault("nano.worker.redis.url", "localhost:6379")
	conf.SetDefault("nano.worker.redis.pool", "3")

	config := config.NewConfig(conf)

	tcp := acceptor.NewTCPAcceptor(fmt.Sprintf(":%d", *port))

	builder := nano.NewBuilderWithConfigs(*isFrontend, *svType, nano.Cluster, map[string]string{}, config)
	if *isFrontend {
		builder.AddAcceptor(tcp)
	}
	app = builder.Build()

	defer app.Shutdown()

	switch *svType {
	case "metagame":
		app.RegisterRemote(&services.Metagame{},
			component.WithName("metagame"),
			component.WithNameFunc(strings.ToLower),
		)
	case "room":
		app.Register(services.NewRoom(app),
			component.WithName("room"),
			component.WithNameFunc(strings.ToLower),
		)
	case "worker":
		configureWorker()
	}

	app.Start()
}
