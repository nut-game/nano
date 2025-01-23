package main

import (
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/nut-game/nano"
	"github.com/nut-game/nano/acceptor"
	"github.com/nut-game/nano/acceptorwrapper"
	"github.com/nut-game/nano/component"
	"github.com/nut-game/nano/config"
	"github.com/nut-game/nano/examples/demo/rate_limiting/services"
	"github.com/nut-game/nano/metrics"
	"github.com/spf13/viper"
)

func createAcceptor(port int, reporters []metrics.Reporter) acceptor.Acceptor {

	// 5 requests in 1 minute. Doesn't make sense, just to test
	// rate limiting
	vConfig := viper.New()
	vConfig.Set("nano.conn.ratelimiting.limit", 5)
	vConfig.Set("nano.conn.ratelimiting.interval", time.Minute)
	pConfig := config.NewConfig(vConfig)

	rateLimitConfig := config.NewNanoConfig(pConfig).Conn.RateLimiting

	tcp := acceptor.NewTCPAcceptor(fmt.Sprintf(":%d", port))
	return acceptorwrapper.WithWrappers(
		tcp,
		acceptorwrapper.NewRateLimitingWrapper(reporters, rateLimitConfig))
}

var app nano.Nano

func main() {
	port := flag.Int("port", 3250, "the port to listen")
	svType := "room"

	flag.Parse()

	config := config.NewDefaultNanoConfig()
	builder := nano.NewDefaultBuilder(true, svType, nano.Cluster, map[string]string{}, *config)
	builder.AddAcceptor(createAcceptor(*port, builder.MetricsReporters))

	app = builder.Build()

	defer app.Shutdown()

	room := services.NewRoom()
	app.Register(room,
		component.WithName("room"),
		component.WithNameFunc(strings.ToLower),
	)

	app.Start()
}
