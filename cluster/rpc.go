package cluster

import (
	"github.com/nut-game/nano/service"
	"github.com/nut-game/nano/session"
)

func (n *Node) NewRpcSession(gateAddr string) (*session.Session, error) {
	sid := service.Connections.SessionID()
	return n.findOrCreateSession(sid, gateAddr)
}
