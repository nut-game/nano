// Copyright (c) nano Author and TFG Co. All Rights Reserved.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package nano

import (
	"context"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"syscall"

	"time"

	"github.com/nut-game/nano/acceptor"
	"github.com/nut-game/nano/cluster"
	"github.com/nut-game/nano/component"
	"github.com/nut-game/nano/config"
	"github.com/nut-game/nano/conn/message"
	"github.com/nut-game/nano/constants"
	pcontext "github.com/nut-game/nano/context"
	"github.com/nut-game/nano/docgenerator"
	"github.com/nut-game/nano/errors"
	"github.com/nut-game/nano/groups"
	"github.com/nut-game/nano/interfaces"
	"github.com/nut-game/nano/logger"
	logging "github.com/nut-game/nano/logger/interfaces"
	"github.com/nut-game/nano/metrics"
	mods "github.com/nut-game/nano/modules"
	"github.com/nut-game/nano/remote"
	"github.com/nut-game/nano/router"
	"github.com/nut-game/nano/serialize"
	"github.com/nut-game/nano/service"
	"github.com/nut-game/nano/session"
	"github.com/nut-game/nano/timer"
	"github.com/nut-game/nano/tracing"
	"github.com/nut-game/nano/worker"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/protobuf/proto"
)

// ServerMode represents a server mode
type ServerMode byte

const (
	_ ServerMode = iota
	// Cluster represents a server running with connection to other servers
	Cluster
	// Standalone represents a server running without connection to other servers
	Standalone
)

// Nano App interface
type Nano interface {
	GetDieChan() chan bool
	SetDebug(debug bool)
	SetHeartbeatTime(interval time.Duration)
	GetServerID() string
	GetMetricsReporters() []metrics.Reporter
	GetServer() *cluster.Server
	GetServerByID(id string) (*cluster.Server, error)
	GetServersByType(t string) (map[string]*cluster.Server, error)
	GetServers() []*cluster.Server
	GetSessionFromCtx(ctx context.Context) session.Session
	Start()
	SetDictionary(dict map[string]uint16) error
	AddRoute(serverType string, routingFunction router.RoutingFunc) error
	Shutdown()
	StartWorker()
	RegisterRPCJob(rpcJob worker.RPCJob) error
	Documentation(getPtrNames bool) (map[string]any, error)
	IsRunning() bool

	RPC(ctx context.Context, routeStr string, reply proto.Message, arg proto.Message) error
	RPCTo(ctx context.Context, serverID, routeStr string, reply proto.Message, arg proto.Message) error
	ReliableRPC(
		routeStr string,
		metadata map[string]any,
		reply, arg proto.Message,
	) (jid string, err error)
	ReliableRPCWithOptions(
		routeStr string,
		metadata map[string]any,
		reply, arg proto.Message,
		opts *config.EnqueueOpts,
	) (jid string, err error)

	SendPushToUsers(route string, v any, uids []string, frontendType string) ([]string, error)
	SendKickToUsers(uids []string, frontendType string) ([]string, error)

	GroupCreate(ctx context.Context, groupName string) error
	GroupCreateWithTTL(ctx context.Context, groupName string, ttlTime time.Duration) error
	GroupMembers(ctx context.Context, groupName string) ([]string, error)
	GroupBroadcast(ctx context.Context, frontendType, groupName, route string, v any) error
	GroupContainsMember(ctx context.Context, groupName, uid string) (bool, error)
	GroupAddMember(ctx context.Context, groupName, uid string) error
	GroupRemoveMember(ctx context.Context, groupName, uid string) error
	GroupRemoveAll(ctx context.Context, groupName string) error
	GroupCountMembers(ctx context.Context, groupName string) (int, error)
	GroupRenewTTL(ctx context.Context, groupName string) error
	GroupDelete(ctx context.Context, groupName string) error

	Register(c component.Component, options ...component.Option)
	RegisterRemote(c component.Component, options ...component.Option)

	RegisterModule(module interfaces.Module, name string) error
	RegisterModuleAfter(module interfaces.Module, name string) error
	RegisterModuleBefore(module interfaces.Module, name string) error
	GetModule(name string) (interfaces.Module, error)

	GetNumberOfConnectedClients() int64
}

// App is the base app struct
type App struct {
	acceptors         []acceptor.Acceptor
	config            config.NanoConfig
	debug             bool
	dieChan           chan bool
	heartbeat         time.Duration
	onSessionBind     func(session.Session)
	router            *router.Router
	rpcClient         cluster.RPCClient
	rpcServer         cluster.RPCServer
	metricsReporters  []metrics.Reporter
	running           bool
	serializer        serialize.Serializer
	server            *cluster.Server
	serverMode        ServerMode
	serviceDiscovery  cluster.ServiceDiscovery
	startAt           time.Time
	worker            *worker.Worker
	remoteService     *service.RemoteService
	handlerService    *service.HandlerService
	handlerComp       []regComp
	remoteComp        []regComp
	modulesMap        map[string]interfaces.Module
	modulesArr        []moduleWrapper
	sessionModulesArr []sessionModuleWrapper
	groups            groups.GroupService
	sessionPool       session.SessionPool
	externalDieChan   chan bool
	sgChan            chan os.Signal
}

// NewApp is the base constructor for a nano app instance
func NewApp(
	serverMode ServerMode,
	serializer serialize.Serializer,
	acceptors []acceptor.Acceptor,
	dieChan chan bool,
	router *router.Router,
	server *cluster.Server,
	rpcClient cluster.RPCClient,
	rpcServer cluster.RPCServer,
	worker *worker.Worker,
	serviceDiscovery cluster.ServiceDiscovery,
	remoteService *service.RemoteService,
	handlerService *service.HandlerService,
	groups groups.GroupService,
	sessionPool session.SessionPool,
	metricsReporters []metrics.Reporter,
	config config.NanoConfig,
) *App {
	app := &App{
		server:            server,
		config:            config,
		rpcClient:         rpcClient,
		rpcServer:         rpcServer,
		worker:            worker,
		serviceDiscovery:  serviceDiscovery,
		remoteService:     remoteService,
		handlerService:    handlerService,
		groups:            groups,
		debug:             false,
		startAt:           time.Now(),
		dieChan:           dieChan,
		acceptors:         acceptors,
		metricsReporters:  metricsReporters,
		serverMode:        serverMode,
		running:           false,
		serializer:        serializer,
		router:            router,
		handlerComp:       make([]regComp, 0),
		remoteComp:        make([]regComp, 0),
		modulesMap:        make(map[string]interfaces.Module),
		modulesArr:        []moduleWrapper{},
		sessionModulesArr: []sessionModuleWrapper{},
		sessionPool:       sessionPool,
		externalDieChan:   make(chan bool),
		sgChan:            make(chan os.Signal, 1),
	}
	if app.heartbeat == time.Duration(0) {
		app.heartbeat = config.Heartbeat.Interval
	}

	app.initSysRemotes()
	return app
}

// GetDieChan gets the channel that the app sinalizes when its going to die.
func (app *App) GetDieChan() chan bool {
	return app.externalDieChan
}

// SetDebug toggles debug on/off
func (app *App) SetDebug(debug bool) {
	app.debug = debug
}

// SetHeartbeatTime sets the heartbeat time
func (app *App) SetHeartbeatTime(interval time.Duration) {
	app.heartbeat = interval
}

// GetServerID returns the generated server id
func (app *App) GetServerID() string {
	return app.server.ID
}

// GetMetricsReporters gets registered metrics reporters
func (app *App) GetMetricsReporters() []metrics.Reporter {
	return app.metricsReporters
}

// GetServer gets the local server instance
func (app *App) GetServer() *cluster.Server {
	return app.server
}

// GetServerByID returns the server with the specified id
func (app *App) GetServerByID(id string) (*cluster.Server, error) {
	return app.serviceDiscovery.GetServer(id)
}

// GetServersByType get all servers of type
func (app *App) GetServersByType(t string) (map[string]*cluster.Server, error) {
	return app.serviceDiscovery.GetServersByType(t)
}

// GetServers get all servers
func (app *App) GetServers() []*cluster.Server {
	return app.serviceDiscovery.GetServers()
}

// IsRunning indicates if the nano app has been initialized. Note: This
// doesn't cover acceptors, only the nano internal registration and modules
// initialization.
func (app *App) IsRunning() bool {
	return app.running
}

// SetLogger logger setter
func SetLogger(l logging.Logger) {
	logger.SetLogger(l)
}

func (app *App) initSysRemotes() {
	sys := remote.NewSys(app.sessionPool)
	app.RegisterRemote(sys,
		component.WithName("sys"),
		component.WithNameFunc(strings.ToLower),
	)
}

func (app *App) periodicMetrics() {
	period := app.config.Metrics.Period
	go metrics.ReportSysMetrics(app.metricsReporters, period)

	if app.worker.Started() {
		go worker.Report(app.metricsReporters, period)
	}
}

// Start starts the app
func (app *App) Start() {
	if !app.server.Frontend && len(app.acceptors) > 0 {
		logger.Fatal("acceptors are not allowed on backend servers")
	}

	if app.server.Frontend && len(app.acceptors) == 0 {
		logger.Fatal("frontend servers should have at least one configured acceptor")
	}

	if app.serverMode == Cluster {
		if reflect.TypeOf(app.rpcClient) == reflect.TypeOf(&cluster.GRPCClient{}) {
			app.serviceDiscovery.AddListener(app.rpcClient.(*cluster.GRPCClient))
		}

		if err := app.RegisterModuleBefore(app.rpcServer, "rpcServer"); err != nil {
			logger.Fatal("failed to register rpc server module: %s", err.Error())
		}
		if err := app.RegisterModuleBefore(app.rpcClient, "rpcClient"); err != nil {
			logger.Fatal("failed to register rpc client module: %s", err.Error())
		}
		// set the service discovery as the last module to be started to ensure
		// all modules have been properly initialized before the server starts
		// receiving requests from other nano servers
		if err := app.RegisterModuleAfter(app.serviceDiscovery, "serviceDiscovery"); err != nil {
			logger.Fatal("failed to register service discovery module: %s", err.Error())
		}
	}

	app.periodicMetrics()

	app.listen()

	defer func() {
		timer.GlobalTicker.Stop()
		app.running = false
	}()

	signal.Notify(app.sgChan, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	maxSessionCount := func() int64 {
		count := app.sessionPool.GetSessionCount()
		mc := app.maxModuleSessionCount()
		if mc > count {
			count = mc
		}
		return count
	}

	// stop server
	select {
	case <-app.dieChan:
		logger.Warn("the app will shutdown in a few seconds")
	case s := <-app.sgChan:
		logger.Warn("got signal: ", s, ", shutting down...")
		if app.config.Session.Drain.Enabled && s == syscall.SIGTERM {
			logger.Info("Session drain is enabled, draining all sessions before shutting down")
			timeoutTimer := time.NewTimer(app.config.Session.Drain.Timeout)
			app.startModuleSessionDraining()
		loop:
			for {
				if maxSessionCount() == 0 {
					logger.Info("All sessions drained")
					break loop
				}
				select {
				case s := <-app.sgChan:
					logger.Warn("got signal: ", s)
					if s == syscall.SIGINT {
						logger.Warnf("Bypassing session draing due to SIGINT. %d sessions will be immediately terminated", maxSessionCount())
					}
					break loop
				case <-timeoutTimer.C:
					logger.Warnf("Session drain has reached maximum timeout. %d sessions will be immediately terminated", maxSessionCount())
					break loop
				case <-time.After(app.config.Session.Drain.Period):
					logger.Infof("Waiting for all sessions to finish: %d sessions remaining...", maxSessionCount())
				}
			}
		}
	}

	app.Shutdown()
	close(app.externalDieChan)
	close(app.sgChan)
	logger.Warn("server is stopping...")

	app.sessionPool.CloseAll()
	app.shutdownModules()
	app.shutdownComponents()
}

func (app *App) listen() {
	app.startupComponents()
	// create global ticker instance, timer precision could be customized
	// by SetTimerPrecision
	timer.GlobalTicker = time.NewTicker(timer.Precision)

	logger.Infof("starting server %s:%s", app.server.Type, app.server.ID)
	for i := 0; i < app.config.Concurrency.Handler.Dispatch; i++ {
		go app.handlerService.Dispatch(i)
	}
	for _, acc := range app.acceptors {
		a := acc
		go func() {
			for conn := range a.GetConnChan() {
				go app.handlerService.Handle(conn)
			}
		}()
		if app.config.Acceptor.ProxyProtocol {
			logger.Info("Enabling PROXY protocol for inbound connections")
			a.EnableProxyProtocol()
		} else {
			logger.Debug("PROXY protocol is disabled for inbound connections")
		}
		go func() {
			a.ListenAndServe()
		}()
		logger.Infof("Waiting for Acceptor %s to start on addr %s", reflect.TypeOf(a), a.GetConfiguredAddress())

		for !a.IsRunning() {
		}

		logger.Infof("Acceptor %s on addr %s is now accepting connections", reflect.TypeOf(a), a.GetAddr())
	}

	if app.serverMode == Cluster && app.server.Frontend && app.config.Session.Unique {
		unique := mods.NewUniqueSession(app.server, app.rpcServer, app.rpcClient, app.sessionPool)
		app.remoteService.AddRemoteBindingListener(unique)
		app.RegisterModule(unique, "uniqueSession")
	}

	app.startModules()

	logger.Info("all modules started!")

	app.running = true
}

// SetDictionary sets routes map
func (app *App) SetDictionary(dict map[string]uint16) error {
	if app.running {
		return constants.ErrChangeDictionaryWhileRunning
	}
	return message.SetDictionary(dict)
}

// AddRoute adds a routing function to a server type
func (app *App) AddRoute(
	serverType string,
	routingFunction router.RoutingFunc,
) error {
	if app.router != nil {
		if app.running {
			return constants.ErrChangeRouteWhileRunning
		}
		app.router.AddRoute(serverType, routingFunction)
	} else {
		return constants.ErrRouterNotInitialized
	}
	return nil
}

// Shutdown send a signal to let 'nano' shutdown itself.
func (app *App) Shutdown() {
	select {
	case <-app.dieChan: // prevent closing closed channel
	default:
		close(app.dieChan)
	}
}

// Error creates a new error with a code, message and metadata
func Error(err error, code string, metadata ...map[string]string) *errors.Error {
	return errors.NewError(err, code, metadata...)
}

// GetSessionFromCtx retrieves a session from a given context
func (app *App) GetSessionFromCtx(ctx context.Context) session.Session {
	sessionVal := ctx.Value(constants.SessionCtxKey)
	if sessionVal == nil {
		logger.Debug("ctx doesn't contain a session, are you calling GetSessionFromCtx from inside a remote?")
		return nil
	}
	return sessionVal.(session.Session)
}

// GetDefaultLoggerFromCtx returns the default logger from the given context
func GetDefaultLoggerFromCtx(ctx context.Context) logging.Logger {
	l := ctx.Value(constants.LoggerCtxKey)
	if l == nil {
		return logger.Log
	}

	return l.(logging.Logger)
}

// AddMetricTagsToPropagateCtx adds a key and metric tags that will
// be propagated through RPC calls. Use the same tags that are at
// 'nano.metrics.additionalLabels' config
func AddMetricTagsToPropagateCtx(
	ctx context.Context,
	tags map[string]string,
) context.Context {
	return pcontext.AddToPropagateCtx(ctx, constants.MetricTagsKey, tags)
}

// AddToPropagateCtx adds a key and value that will be propagated through RPC calls
func AddToPropagateCtx(ctx context.Context, key string, val any) context.Context {
	return pcontext.AddToPropagateCtx(ctx, key, val)
}

// GetFromPropagateCtx adds a key and value that came through RPC calls
func GetFromPropagateCtx(ctx context.Context, key string) any {
	return pcontext.GetFromPropagateCtx(ctx, key)
}

// ExtractSpan retrieves an opentracing span context from the given context
// The span context can be received directly or via an RPC call
func ExtractSpan(ctx context.Context) (opentracing.SpanContext, error) {
	return tracing.ExtractSpan(ctx)
}

// Documentation returns handler and remotes documentacion
func (app *App) Documentation(getPtrNames bool) (map[string]any, error) {
	handlerDocs, err := app.handlerService.Docs(getPtrNames)
	if err != nil {
		return nil, err
	}
	remoteDocs, err := app.remoteService.Docs(getPtrNames)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"handlers": handlerDocs,
		"remotes":  remoteDocs,
	}, nil
}

// AddGRPCInfoToMetadata adds host, external host and
// port into metadata
func AddGRPCInfoToMetadata(
	metadata map[string]string,
	region string,
	host, port string,
	externalHost, externalPort string,
) map[string]string {
	metadata[constants.GRPCHostKey] = host
	metadata[constants.GRPCPortKey] = port
	metadata[constants.GRPCExternalHostKey] = externalHost
	metadata[constants.GRPCExternalPortKey] = externalPort
	metadata[constants.RegionKey] = region
	return metadata
}

// Descriptor returns the protobuf message descriptor for a given message name
func Descriptor(protoName string) ([]byte, error) {
	return docgenerator.ProtoDescriptors(protoName)
}

// StartWorker configures, starts and returns nano worker
func (app *App) StartWorker() {
	app.worker.Start()
}

// RegisterRPCJob registers rpc job to execute jobs with retries
func (app *App) RegisterRPCJob(rpcJob worker.RPCJob) error {
	err := app.worker.RegisterRPCJob(rpcJob)
	return err
}

// GetNumberOfConnectedClients returns the number of connected clients
func (app *App) GetNumberOfConnectedClients() int64 {
	return app.sessionPool.GetSessionCount()
}
