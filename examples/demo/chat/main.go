package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"strings"

	"github.com/nut-game/nano"
	"github.com/nut-game/nano/acceptor"
	"github.com/nut-game/nano/component"
	"github.com/nut-game/nano/config"
	"github.com/nut-game/nano/groups"
	"github.com/nut-game/nano/logger"
	"github.com/nut-game/nano/timer"
)

type (
	// Room represents a component that contains a bundle of room related handler
	// like Join/Message
	Room struct {
		component.Base
		timer *timer.Timer
		app   nano.Nano
	}

	// UserMessage represents a message that user sent
	UserMessage struct {
		Name    string `json:"name"`
		Content string `json:"content"`
	}

	// NewUser message will be received when new user join room
	NewUser struct {
		Content string `json:"content"`
	}

	// AllMembers contains all members uid
	AllMembers struct {
		Members []string `json:"members"`
	}

	// JoinResponse represents the result of joining room
	JoinResponse struct {
		Code   int    `json:"code"`
		Result string `json:"result"`
	}
)

// NewRoom returns a Handler Base implementation
func NewRoom(app nano.Nano) *Room {
	return &Room{
		app: app,
	}
}

// AfterInit component lifetime callback
func (r *Room) AfterInit() {
	r.timer = nano.NewTimer(time.Minute, func() {
		count, err := r.app.GroupCountMembers(context.Background(), "room")
		logger.Debugf("UserCount: Time=> %s, Count=> %d, Error=> %v", time.Now().String(), count, err)
	})
}

// Join room
func (r *Room) Join(ctx context.Context, msg []byte) (*JoinResponse, error) {
	s := r.app.GetSessionFromCtx(ctx)
	fakeUID := s.ID()                              // just use s.ID as uid !!!
	err := s.Bind(ctx, strconv.Itoa(int(fakeUID))) // binding session uid

	if err != nil {
		return nil, nano.Error(err, "RH-000", map[string]string{"failed": "bind"})
	}

	uids, err := r.app.GroupMembers(ctx, "room")
	if err != nil {
		return nil, err
	}
	s.Push("onMembers", &AllMembers{Members: uids})
	// notify others
	r.app.GroupBroadcast(ctx, "chat", "room", "onNewUser", &NewUser{Content: fmt.Sprintf("New user: %s", s.UID())})
	// new user join group
	r.app.GroupAddMember(ctx, "room", s.UID()) // add session to group

	// on session close, remove it from group
	s.OnClose(func() {
		r.app.GroupRemoveMember(ctx, "room", s.UID())
	})

	return &JoinResponse{Result: "success"}, nil
}

// Message sync last message to all members
func (r *Room) Message(ctx context.Context, msg *UserMessage) {
	err := r.app.GroupBroadcast(ctx, "chat", "room", "onMessage", msg)
	if err != nil {
		fmt.Println("error broadcasting message", err)
	}
}

var app nano.Nano

func main() {
	conf := configApp()
	builder := nano.NewDefaultBuilder(true, "chat", nano.Cluster, map[string]string{}, *conf)
	builder.AddAcceptor(acceptor.NewWSAcceptor(":40001"))
	builder.Groups = groups.NewMemoryGroupService(builder.Config.Groups.Memory)
	app = builder.Build()

	defer app.Shutdown()

	err := app.GroupCreate(context.Background(), "room")
	if err != nil {
		panic(err)
	}

	// rewrite component and handler name
	room := NewRoom(app)
	app.Register(room,
		component.WithName("room"),
		component.WithNameFunc(strings.ToLower),
	)

	log.SetFlags(log.LstdFlags | log.Llongfile)

	http.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("web"))))

	go http.ListenAndServe(":4000", nil)

	app.Start()
}

func configApp() *config.NanoConfig {
	conf := config.NewDefaultNanoConfig()
	conf.Buffer.Handler.LocalProcess = 15
	conf.Heartbeat.Interval = time.Duration(15 * time.Second)
	conf.Buffer.Agent.Messages = 32
	conf.Handler.Messages.Compression = false
	return conf
}
