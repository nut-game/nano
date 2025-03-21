package services

import (
	"context"
	"fmt"

	"github.com/nut-game/nano"
	"github.com/nut-game/nano/component"
	"github.com/nut-game/nano/examples/demo/protos"
)

// ConnectorRemote is a remote that will receive rpc's
type ConnectorRemote struct {
	component.Base
}

// Connector struct
type Connector struct {
	component.Base
	app nano.Nano
}

// SessionData struct
type SessionData struct {
	Data map[string]any
}

// Response struct
type Response struct {
	Code int32
	Msg  string
}

// NewConnector ctor
func NewConnector(app nano.Nano) *Connector {
	return &Connector{app: app}
}

func reply(code int32, msg string) (*Response, error) {
	res := &Response{
		Code: code,
		Msg:  msg,
	}
	return res, nil
}

// GetSessionData gets the session data
func (c *Connector) GetSessionData(ctx context.Context) (*SessionData, error) {
	s := c.app.GetSessionFromCtx(ctx)
	res := &SessionData{
		Data: s.GetData(),
	}
	return res, nil
}

// SetSessionData sets the session data
func (c *Connector) SetSessionData(ctx context.Context, data *SessionData) (*Response, error) {
	s := c.app.GetSessionFromCtx(ctx)
	err := s.SetData(data.Data)
	if err != nil {
		return nil, nano.Error(err, "CN-000", map[string]string{"failed": "set data"})
	}
	return reply(200, "success")
}

// NotifySessionData sets the session data
func (c *Connector) NotifySessionData(ctx context.Context, data *SessionData) {
	s := c.app.GetSessionFromCtx(ctx)
	err := s.SetData(data.Data)
	if err != nil {
		fmt.Println("got error on notify", err)
	}
}

// SendPushToUser sends a push to a user
func (c *Connector) SendPushToUser(ctx context.Context, msg *UserMessage) (*Response, error) {
	_, err := c.app.SendPushToUsers("onMessage", msg, []string{"2"}, "connector")
	if err != nil {
		return nil, err
	}
	return &Response{
		Code: 200,
		Msg:  "boa",
	}, nil
}

// RemoteFunc is a function that will be called remotely
func (c *ConnectorRemote) RemoteFunc(ctx context.Context, msg *protos.RPCMsg) (*protos.RPCRes, error) {
	fmt.Printf("received a remote call with this message: %s\n", msg)
	return &protos.RPCRes{
		Msg: fmt.Sprintf("received msg: %s", msg.GetMsg()),
	}, nil
}
