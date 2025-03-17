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
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func setup() {
	os.Setenv("NANO_HEARTBEAT_INTERVAL", "123s")
	os.Setenv("NANO_CONCURRENCY_TEST", "42")
	os.Setenv("NANO_BUFFER_TEST", "14")
}

func shutdown() {
	os.Unsetenv("NANO_HEARTBEAT_INTERVAL")
	os.Unsetenv("NANO_CONCURRENCY_TEST")
	os.Unsetenv("NANO_BUFFER_TEST")
}

func TestNewConfig(t *testing.T) {
	t.Parallel()

	cfg := viper.New()
	cfg.SetDefault("nano.buffer.agent.messages", 20)
	cfg.Set("nano.concurrency.handler.dispatch", 23)
	cfg.SetDefault("nano.no.default", "custom")

	tables := []struct {
		in  []*viper.Viper
		key string
		val any
	}{
		{[]*viper.Viper{}, "nano.buffer.agent.messages", 100},
		{[]*viper.Viper{cfg}, "nano.buffer.agent.messages", 20},
		{[]*viper.Viper{}, "nano.no.default", nil},
		{[]*viper.Viper{cfg}, "nano.no.default", "custom"},
		{[]*viper.Viper{}, "nano.concurrency.handler.dispatch", 25},
		{[]*viper.Viper{cfg}, "nano.concurrency.handler.dispatch", 23},
		{[]*viper.Viper{}, "nano.heartbeat.interval", "123s"},
		{[]*viper.Viper{cfg}, "nano.heartbeat.interval", "123s"},
		{[]*viper.Viper{}, "nano.concurrency.test", "42"},
		{[]*viper.Viper{cfg}, "nano.concurrency.test", "42"},
		{[]*viper.Viper{}, "nano.buffer.test", "14"},
		{[]*viper.Viper{cfg}, "nano.buffer.test", "14"},
	}

	for _, table := range tables {
		t.Run(fmt.Sprintf("arguments:%d", len(table.in)), func(t *testing.T) {
			c := NewConfig(table.in...)
			assert.Equal(t, table.val, c.Get(table.key))
		})
	}
}

func TestGetDuration(t *testing.T) {
	t.Parallel()

	c := NewConfig()
	tables := []struct {
		key string
		val time.Duration
	}{
		{"nano.heartbeat.interval", 123 * time.Second},
		{"nano.cluster.sd.etcd.dialtimeout", 5 * time.Second},
		{"unexistent", time.Duration(0)},
	}

	for _, table := range tables {
		t.Run(fmt.Sprintf("key:%s val:%d", table.key, table.val), func(t *testing.T) {
			assert.Equal(t, table.val, c.GetDuration(table.key))
		})
	}
}

func TestGetString(t *testing.T) {
	t.Parallel()

	c := NewConfig()
	tables := []struct {
		key string
		val []string
	}{
		{"nano.cluster.sd.etcd.endpoints", []string{"localhost:2379"}},
		{"unexistent", nil},
	}

	for _, table := range tables {
		t.Run(fmt.Sprintf("key:%s val:%s", table.key, table.val), func(t *testing.T) {
			assert.Equal(t, table.val, c.GetStringSlice(table.key))
		})
	}
}

func TestGetInt(t *testing.T) {
	t.Parallel()

	c := NewConfig()
	tables := []struct {
		key string
		val int
	}{
		{"nano.buffer.agent.messages", 100},
		{"unexistent", 0},
	}

	for _, table := range tables {
		t.Run(fmt.Sprintf("key:%s val:%d", table.key, table.val), func(t *testing.T) {
			assert.Equal(t, table.val, c.GetInt(table.key))
		})
	}
}

func TestGetStringSlice(t *testing.T) {
	t.Parallel()

	c := NewConfig()
	tables := []struct {
		key string
		val []string
	}{
		{"nano.cluster.sd.etcd.endpoints", []string{"localhost:2379"}},
		{"unexistent", nil},
	}

	for _, table := range tables {
		t.Run(fmt.Sprintf("key:%s val:%s", table.key, table.val), func(t *testing.T) {
			assert.Equal(t, table.val, c.GetStringSlice(table.key))
		})
	}
}

func TestGet(t *testing.T) {
	t.Parallel()

	c := NewConfig()
	tables := []struct {
		key string
		val any
	}{
		{"nano.buffer.agent.messages", 100},
		{"unexistent", nil},
	}

	for _, table := range tables {
		t.Run(fmt.Sprintf("key:%s val:%v", table.key, table.val), func(t *testing.T) {
			assert.Equal(t, table.val, c.Get(table.key))
		})
	}
}
