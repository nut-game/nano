// Copyright (c) TFG Co. All Rights Reserved.
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

package config

import (
	"reflect"
	"strings"

	"github.com/mitchellh/mapstructure"

	"github.com/spf13/viper"
)

// Config is a wrapper around a viper config
type Config struct {
	viper.Viper
}

// NewConfig creates a new config with a given viper config if given
func NewConfig(cfgs ...*viper.Viper) *Config {
	var cfg *viper.Viper
	if len(cfgs) > 0 {
		cfg = cfgs[0]
	} else {
		cfg = viper.New()
	}

	cfg.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	cfg.AutomaticEnv()
	c := &Config{*cfg}
	c.fillDefaultValues()
	return c
}

func (c *Config) fillDefaultValues() {
	config := NewDefaultNanoConfig()

	defaultsMap := map[string]interface{}{
		"nano.serializertype":        config.SerializerType,
		"nano.buffer.agent.messages": config.Buffer.Agent.Messages,
		// the max buffer size that nats will accept, if this buffer overflows, messages will begin to be dropped
		"nano.buffer.handler.localprocess":                    config.Buffer.Handler.LocalProcess,
		"nano.buffer.handler.remoteprocess":                   config.Buffer.Handler.RemoteProcess,
		"nano.cluster.info.region":                            config.Cluster.Info.Region,
		"nano.cluster.rpc.client.grpc.dialtimeout":            config.Cluster.RPC.Client.Grpc.DialTimeout,
		"nano.cluster.rpc.client.grpc.requesttimeout":         config.Cluster.RPC.Client.Grpc.RequestTimeout,
		"nano.cluster.rpc.client.grpc.lazyconnection":         config.Cluster.RPC.Client.Grpc.LazyConnection,
		"nano.cluster.rpc.client.nats.connect":                config.Cluster.RPC.Client.Nats.Connect,
		"nano.cluster.rpc.client.nats.connectiontimeout":      config.Cluster.RPC.Client.Nats.ConnectionTimeout,
		"nano.cluster.rpc.client.nats.maxreconnectionretries": config.Cluster.RPC.Client.Nats.MaxReconnectionRetries,
		"nano.cluster.rpc.client.nats.websocketcompression":   config.Cluster.RPC.Client.Nats.WebsocketCompression,
		"nano.cluster.rpc.client.nats.reconnectjitter":        config.Cluster.RPC.Client.Nats.ReconnectJitter,
		"nano.cluster.rpc.client.nats.reconnectjittertls":     config.Cluster.RPC.Client.Nats.ReconnectJitterTLS,
		"nano.cluster.rpc.client.nats.reconnectwait":          config.Cluster.RPC.Client.Nats.ReconnectWait,
		"nano.cluster.rpc.client.nats.pinginterval":           config.Cluster.RPC.Client.Nats.PingInterval,
		"nano.cluster.rpc.client.nats.maxpingsoutstanding":    config.Cluster.RPC.Client.Nats.MaxPingsOutstanding,
		"nano.cluster.rpc.client.nats.requesttimeout":         config.Cluster.RPC.Client.Nats.RequestTimeout,
		"nano.cluster.rpc.server.grpc.port":                   config.Cluster.RPC.Server.Grpc.Port,
		"nano.cluster.rpc.server.nats.connect":                config.Cluster.RPC.Server.Nats.Connect,
		"nano.cluster.rpc.server.nats.connectiontimeout":      config.Cluster.RPC.Server.Nats.ConnectionTimeout,
		"nano.cluster.rpc.server.nats.maxreconnectionretries": config.Cluster.RPC.Server.Nats.MaxReconnectionRetries,
		"nano.cluster.rpc.server.nats.websocketcompression":   config.Cluster.RPC.Server.Nats.WebsocketCompression,
		"nano.cluster.rpc.server.nats.reconnectjitter":        config.Cluster.RPC.Server.Nats.ReconnectJitter,
		"nano.cluster.rpc.server.nats.reconnectjittertls":     config.Cluster.RPC.Server.Nats.ReconnectJitterTLS,
		"nano.cluster.rpc.server.nats.reconnectwait":          config.Cluster.RPC.Server.Nats.ReconnectWait,
		"nano.cluster.rpc.server.nats.pinginterval":           config.Cluster.RPC.Server.Nats.PingInterval,
		"nano.cluster.rpc.server.nats.maxpingsoutstanding":    config.Cluster.RPC.Server.Nats.MaxPingsOutstanding,
		"nano.cluster.rpc.server.nats.services":               config.Cluster.RPC.Server.Nats.Services,
		"nano.cluster.rpc.server.nats.buffer.messages":        config.Cluster.RPC.Server.Nats.Buffer.Messages,
		"nano.cluster.rpc.server.nats.buffer.push":            config.Cluster.RPC.Server.Nats.Buffer.Push,
		"nano.cluster.sd.etcd.dialtimeout":                    config.Cluster.SD.Etcd.DialTimeout,
		"nano.cluster.sd.etcd.endpoints":                      config.Cluster.SD.Etcd.Endpoints,
		"nano.cluster.sd.etcd.prefix":                         config.Cluster.SD.Etcd.Prefix,
		"nano.cluster.sd.etcd.grantlease.maxretries":          config.Cluster.SD.Etcd.GrantLease.MaxRetries,
		"nano.cluster.sd.etcd.grantlease.retryinterval":       config.Cluster.SD.Etcd.GrantLease.RetryInterval,
		"nano.cluster.sd.etcd.grantlease.timeout":             config.Cluster.SD.Etcd.GrantLease.Timeout,
		"nano.cluster.sd.etcd.heartbeat.log":                  config.Cluster.SD.Etcd.Heartbeat.Log,
		"nano.cluster.sd.etcd.heartbeat.ttl":                  config.Cluster.SD.Etcd.Heartbeat.TTL,
		"nano.cluster.sd.etcd.revoke.timeout":                 config.Cluster.SD.Etcd.Revoke.Timeout,
		"nano.cluster.sd.etcd.syncservers.interval":           config.Cluster.SD.Etcd.SyncServers.Interval,
		"nano.cluster.sd.etcd.syncservers.parallelism":        config.Cluster.SD.Etcd.SyncServers.Parallelism,
		"nano.cluster.sd.etcd.shutdown.delay":                 config.Cluster.SD.Etcd.Shutdown.Delay,
		"nano.cluster.sd.etcd.servertypeblacklist":            config.Cluster.SD.Etcd.ServerTypesBlacklist,
		// the sum of this config among all the frontend servers should always be less than
		// the sum of nano.buffer.cluster.rpc.server.nats.messages, for covering the worst case scenario
		// a single backend server should have the config nano.buffer.cluster.rpc.server.nats.messages bigger
		// than the sum of the config nano.concurrency.handler.dispatch among all frontend servers
		"nano.acceptor.proxyprotocol":                    config.Acceptor.ProxyProtocol,
		"nano.concurrency.handler.dispatch":              config.Concurrency.Handler.Dispatch,
		"nano.defaultpipelines.structvalidation.enabled": config.DefaultPipelines.StructValidation.Enabled,
		"nano.groups.etcd.dialtimeout":                   config.Groups.Etcd.DialTimeout,
		"nano.groups.etcd.endpoints":                     config.Groups.Etcd.Endpoints,
		"nano.groups.etcd.prefix":                        config.Groups.Etcd.Prefix,
		"nano.groups.etcd.transactiontimeout":            config.Groups.Etcd.TransactionTimeout,
		"nano.groups.memory.tickduration":                config.Groups.Memory.TickDuration,
		"nano.handler.messages.compression":              config.Handler.Messages.Compression,
		"nano.heartbeat.interval":                        config.Heartbeat.Interval,
		"nano.metrics.additionalLabels":                  config.Metrics.AdditionalLabels,
		"nano.metrics.constLabels":                       config.Metrics.ConstLabels,
		"nano.metrics.custom":                            config.Metrics.Custom,
		"nano.metrics.period":                            config.Metrics.Period,
		"nano.metrics.prometheus.enabled":                config.Metrics.Prometheus.Enabled,
		"nano.metrics.prometheus.port":                   config.Metrics.Prometheus.Port,
		"nano.metrics.statsd.enabled":                    config.Metrics.Statsd.Enabled,
		"nano.metrics.statsd.host":                       config.Metrics.Statsd.Host,
		"nano.metrics.statsd.prefix":                     config.Metrics.Statsd.Prefix,
		"nano.metrics.statsd.rate":                       config.Metrics.Statsd.Rate,
		"nano.modules.bindingstorage.etcd.dialtimeout":   config.Modules.BindingStorage.Etcd.DialTimeout,
		"nano.modules.bindingstorage.etcd.endpoints":     config.Modules.BindingStorage.Etcd.Endpoints,
		"nano.modules.bindingstorage.etcd.leasettl":      config.Modules.BindingStorage.Etcd.LeaseTTL,
		"nano.modules.bindingstorage.etcd.prefix":        config.Modules.BindingStorage.Etcd.Prefix,
		"nano.conn.ratelimiting.limit":                   config.Conn.RateLimiting.Limit,
		"nano.conn.ratelimiting.interval":                config.Conn.RateLimiting.Interval,
		"nano.conn.ratelimiting.forcedisable":            config.Conn.RateLimiting.ForceDisable,
		"nano.session.unique":                            config.Session.Unique,
		"nano.session.drain.enabled":                     config.Session.Drain.Enabled,
		"nano.session.drain.timeout":                     config.Session.Drain.Timeout,
		"nano.session.drain.period":                      config.Session.Drain.Period,
		"nano.worker.concurrency":                        config.Worker.Concurrency,
		"nano.worker.redis.pool":                         config.Worker.Redis.Pool,
		"nano.worker.redis.url":                          config.Worker.Redis.ServerURL,
		"nano.worker.retry.enabled":                      config.Worker.Retry.Enabled,
		"nano.worker.retry.exponential":                  config.Worker.Retry.Exponential,
		"nano.worker.retry.max":                          config.Worker.Retry.Max,
		"nano.worker.retry.maxDelay":                     config.Worker.Retry.MaxDelay,
		"nano.worker.retry.maxRandom":                    config.Worker.Retry.MaxRandom,
		"nano.worker.retry.minDelay":                     config.Worker.Retry.MinDelay,
	}

	for param := range defaultsMap {
		val := c.Get(param)
		if val == nil {
			c.SetDefault(param, defaultsMap[param])
		} else {
			c.SetDefault(param, val)
			c.Set(param, val)
		}

	}
}

// UnmarshalKey unmarshals key into v
func (c *Config) UnmarshalKey(key string, rawVal interface{}) error {
	key = strings.ToLower(key)
	delimiter := "."
	prefix := key + delimiter

	i := c.Get(key)
	if i == nil {
		return nil
	}
	if isStringMapInterface(i) {
		val := i.(map[string]interface{})
		keys := c.AllKeys()
		for _, k := range keys {
			if !strings.HasPrefix(k, prefix) {
				continue
			}
			mk := strings.TrimPrefix(k, prefix)
			mk = strings.Split(mk, delimiter)[0]
			if _, exists := val[mk]; exists {
				continue
			}
			mv := c.Get(key + delimiter + mk)
			if mv == nil {
				continue
			}
			val[mk] = mv
		}
		i = val
	}
	return decode(i, defaultDecoderConfig(rawVal))
}

func isStringMapInterface(val interface{}) bool {
	vt := reflect.TypeOf(val)
	return vt.Kind() == reflect.Map &&
		vt.Key().Kind() == reflect.String &&
		vt.Elem().Kind() == reflect.Interface
}

// A wrapper around mapstructure.Decode that mimics the WeakDecode functionality
func decode(input interface{}, config *mapstructure.DecoderConfig) error {
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}
	return decoder.Decode(input)
}

// defaultDecoderConfig returns default mapstructure.DecoderConfig with support
// of time.Duration values & string slices
func defaultDecoderConfig(output interface{}, opts ...viper.DecoderConfigOption) *mapstructure.DecoderConfig {
	c := &mapstructure.DecoderConfig{
		Metadata:         nil,
		Result:           output,
		WeaklyTypedInput: true,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToSliceHookFunc(","),
		),
	}
	for _, opt := range opts {
		opt(c)
	}
	return c

}
