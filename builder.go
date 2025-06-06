package nano

import (
	"github.com/google/uuid"
	"github.com/nut-game/nano/acceptor"
	"github.com/nut-game/nano/agent"
	"github.com/nut-game/nano/cluster"
	"github.com/nut-game/nano/config"
	"github.com/nut-game/nano/conn/codec"
	"github.com/nut-game/nano/conn/message"
	"github.com/nut-game/nano/defaultpipelines"
	"github.com/nut-game/nano/groups"
	"github.com/nut-game/nano/logger"
	"github.com/nut-game/nano/metrics"
	"github.com/nut-game/nano/metrics/models"
	"github.com/nut-game/nano/pipeline"
	"github.com/nut-game/nano/router"
	"github.com/nut-game/nano/serialize"
	"github.com/nut-game/nano/service"
	"github.com/nut-game/nano/session"
	"github.com/nut-game/nano/worker"
)

// Builder holds dependency instances for a nano App
type Builder struct {
	acceptors        []acceptor.Acceptor
	postBuildHooks   []func(app Nano)
	Config           config.NanoConfig
	DieChan          chan bool
	PacketDecoder    codec.PacketDecoder
	PacketEncoder    codec.PacketEncoder
	MessageEncoder   *message.MessagesEncoder
	Serializer       serialize.Serializer
	Router           *router.Router
	RPCClient        cluster.RPCClient
	RPCServer        cluster.RPCServer
	MetricsReporters []metrics.Reporter
	Server           *cluster.Server
	ServerMode       ServerMode
	ServiceDiscovery cluster.ServiceDiscovery
	Groups           groups.GroupService
	SessionPool      session.SessionPool
	Worker           *worker.Worker
	RemoteHooks      *pipeline.RemoteHooks
	HandlerHooks     *pipeline.HandlerHooks
}

// NanoBuilder Builder interface
type NanoBuilder interface {
	// AddPostBuildHook adds a post-build hook to the builder, a function receiving a Nano instance as parameter.
	AddPostBuildHook(hook func(app Nano))
	Build() Nano
}

// NewBuilderWithConfigs return a builder instance with default dependency instances for a nano App
// with configs defined by a config file (config.Config) and default paths (see documentation).
func NewBuilderWithConfigs(
	isFrontend bool,
	serverType string,
	serverMode ServerMode,
	serverMetadata map[string]string,
	conf *config.Config,
) *Builder {
	config := config.NewNanoConfig(conf)
	return NewBuilder(
		isFrontend,
		serverType,
		serverMode,
		serverMetadata,
		*config,
	)
}

// NewDefaultBuilder return a builder instance with default dependency instances for a nano App,
// with default configs
func NewDefaultBuilder(isFrontend bool, serverType string, serverMode ServerMode, serverMetadata map[string]string, nanoConfig config.NanoConfig) *Builder {
	return NewBuilder(
		isFrontend,
		serverType,
		serverMode,
		serverMetadata,
		nanoConfig,
	)
}

// NewBuilder return a builder instance with default dependency instances for a nano App,
// with configs explicitly defined
func NewBuilder(isFrontend bool,
	serverType string,
	serverMode ServerMode,
	serverMetadata map[string]string,
	config config.NanoConfig,
) *Builder {
	server := cluster.NewServerWithOptions(
		uuid.New().String(),
		serverType,
		isFrontend,
		cluster.WithMetadata(serverMetadata),
		cluster.WithLoopbackEnabled(config.Cluster.RPC.Server.LoopbackEnabled),
	)

	dieChan := make(chan bool)

	metricsReporters := []metrics.Reporter{}
	if config.Metrics.Prometheus.Enabled {
		metricsReporters = addDefaultPrometheus(config.Metrics, config.Metrics.Custom, metricsReporters, serverType)
	}

	if config.Metrics.Statsd.Enabled {
		metricsReporters = addDefaultStatsd(config.Metrics, metricsReporters, serverType)
	}

	handlerHooks := pipeline.NewHandlerHooks()
	if config.DefaultPipelines.StructValidation.Enabled {
		configureDefaultPipelines(handlerHooks)
	}

	sessionPool := session.NewSessionPool()

	var serviceDiscovery cluster.ServiceDiscovery
	var rpcServer cluster.RPCServer
	var rpcClient cluster.RPCClient
	if serverMode == Cluster {
		var err error
		serviceDiscovery, err = cluster.NewEtcdServiceDiscovery(config.Cluster.SD.Etcd, server, dieChan)
		if err != nil {
			logger.Fatalf("error creating default cluster service discovery component: %s", err.Error())
		}

		rpcServer, err = cluster.NewNatsRPCServer(config.Cluster.RPC.Server.Nats, server, metricsReporters, dieChan, sessionPool)
		if err != nil {
			logger.Fatalf("error setting default cluster rpc server component: %s", err.Error())
		}

		rpcClient, err = cluster.NewNatsRPCClient(config.Cluster.RPC.Client.Nats, server, metricsReporters, dieChan)
		if err != nil {
			logger.Fatalf("error setting default cluster rpc client component: %s", err.Error())
		}
	}

	worker, err := worker.NewWorker(config.Worker, config.Worker.Retry)
	if err != nil {
		logger.Fatalf("error creating default worker: %s", err.Error())
	}

	gsi := groups.NewMemoryGroupService(config.Groups.Memory)
	if err != nil {
		panic(err)
	}

	serializer, err := serialize.NewSerializer(serialize.Type(config.SerializerType))
	if err != nil {
		logger.Fatalf("error creating serializer: %s", err.Error())
	}

	return &Builder{
		acceptors:        []acceptor.Acceptor{},
		postBuildHooks:   make([]func(app Nano), 0),
		Config:           config,
		DieChan:          dieChan,
		PacketDecoder:    codec.NewPomeloPacketDecoder(),
		PacketEncoder:    codec.NewPomeloPacketEncoder(),
		MessageEncoder:   message.NewMessagesEncoder(config.Handler.Messages.Compression),
		Serializer:       serializer,
		Router:           router.New(),
		RPCClient:        rpcClient,
		RPCServer:        rpcServer,
		MetricsReporters: metricsReporters,
		Server:           server,
		ServerMode:       serverMode,
		Groups:           gsi,
		RemoteHooks:      pipeline.NewRemoteHooks(),
		HandlerHooks:     handlerHooks,
		ServiceDiscovery: serviceDiscovery,
		SessionPool:      sessionPool,
		Worker:           worker,
	}
}

// AddAcceptor adds a new acceptor to app
func (builder *Builder) AddAcceptor(ac acceptor.Acceptor) {
	if !builder.Server.Frontend {
		logger.Error("tried to add an acceptor to a backend server, skipping")
		return
	}
	builder.acceptors = append(builder.acceptors, ac)
}

// AddPostBuildHook adds a post-build hook to the builder, a function receiving a Nano instance as parameter.
func (builder *Builder) AddPostBuildHook(hook func(app Nano)) {
	builder.postBuildHooks = append(builder.postBuildHooks, hook)
}

// Build returns a valid App instance
func (builder *Builder) Build() Nano {
	handlerPool := service.NewHandlerPool()
	var remoteService *service.RemoteService
	if builder.ServerMode == Standalone {
		if builder.ServiceDiscovery != nil || builder.RPCClient != nil || builder.RPCServer != nil {
			panic("Standalone mode can't have RPC or service discovery instances")
		}
	} else {
		if !(builder.ServiceDiscovery != nil && builder.RPCClient != nil && builder.RPCServer != nil) {
			panic("Cluster mode must have RPC and service discovery instances")
		}

		builder.Router.SetServiceDiscovery(builder.ServiceDiscovery)

		remoteService = service.NewRemoteService(
			builder.RPCClient,
			builder.RPCServer,
			builder.ServiceDiscovery,
			builder.PacketEncoder,
			builder.Serializer,
			builder.Router,
			builder.MessageEncoder,
			builder.Server,
			builder.SessionPool,
			builder.RemoteHooks,
			builder.HandlerHooks,
			handlerPool,
		)

		builder.RPCServer.SetNanoServer(remoteService)
	}

	agentFactory := agent.NewAgentFactory(builder.DieChan,
		builder.PacketDecoder,
		builder.PacketEncoder,
		builder.Serializer,
		builder.Config.Heartbeat.Interval,
		builder.Config.Buffer.Agent.WriteTimeout,
		builder.MessageEncoder,
		builder.Config.Buffer.Agent.Messages,
		builder.SessionPool,
		builder.MetricsReporters,
	)

	handlerService := service.NewHandlerService(
		builder.PacketDecoder,
		builder.Serializer,
		builder.Config.Buffer.Handler.LocalProcess,
		builder.Config.Buffer.Handler.RemoteProcess,
		builder.Server,
		remoteService,
		agentFactory,
		builder.MetricsReporters,
		builder.HandlerHooks,
		handlerPool,
	)

	app := NewApp(
		builder.ServerMode,
		builder.Serializer,
		builder.acceptors,
		builder.DieChan,
		builder.Router,
		builder.Server,
		builder.RPCClient,
		builder.RPCServer,
		builder.Worker,
		builder.ServiceDiscovery,
		remoteService,
		handlerService,
		builder.Groups,
		builder.SessionPool,
		builder.MetricsReporters,
		builder.Config,
	)

	for _, postBuildHook := range builder.postBuildHooks {
		postBuildHook(app)
	}

	return app
}

// NewDefaultApp returns a default nano app instance
func NewDefaultApp(isFrontend bool, serverType string, serverMode ServerMode, serverMetadata map[string]string, config config.NanoConfig) Nano {
	builder := NewDefaultBuilder(isFrontend, serverType, serverMode, serverMetadata, config)
	return builder.Build()
}

func configureDefaultPipelines(handlerHooks *pipeline.HandlerHooks) {
	handlerHooks.BeforeHandler.PushBack(defaultpipelines.StructValidatorInstance.Validate)
}

func addDefaultPrometheus(config config.MetricsConfig, customMetrics models.CustomMetricsSpec, reporters []metrics.Reporter, serverType string) []metrics.Reporter {
	prometheus, err := CreatePrometheusReporter(serverType, config, customMetrics)
	if err != nil {
		logger.Errorf("failed to start prometheus metrics reporter, skipping %v", err)
	} else {
		reporters = append(reporters, prometheus)
	}
	return reporters
}

func addDefaultStatsd(config config.MetricsConfig, reporters []metrics.Reporter, serverType string) []metrics.Reporter {
	statsd, err := CreateStatsdReporter(serverType, config)
	if err != nil {
		logger.Errorf("failed to start statsd metrics reporter, skipping %v", err)
	} else {
		reporters = append(reporters, statsd)
	}
	return reporters
}
