package client

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/nut-game/nano/conn/message"
	"github.com/nut-game/nano/helpers"
	"github.com/nut-game/nano/mocks"
	"github.com/stretchr/testify/assert"
)

func TestSendRequestShouldTimeout(t *testing.T) {
	c := New(100 * time.Millisecond)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConn := mocks.NewMockPlayerConn(ctrl)
	c.conn = mockConn
	go c.pendingRequestsReaper()

	route := "com.sometest.route"
	data := []byte{0x02, 0x03, 0x04}

	m := message.Message{
		Type:  message.Request,
		ID:    1,
		Route: route,
		Data:  data,
		Err:   false,
	}

	pkt, err := c.buildPacket(m)
	assert.NoError(t, err)

	mockConn.EXPECT().Write(pkt)

	c.IncomingMsgChan = make(chan *message.Message, 10)

	c.nextID = 0
	c.SendRequest(route, data)

	msg := helpers.ShouldEventuallyReceive(t, c.IncomingMsgChan, 2*time.Second).(*message.Message)

	assert.Equal(t, true, msg.Err)
}
