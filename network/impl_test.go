package network

import (
	"reflect"
	"testing"

	"github.com/shoooooman/mg-rs/common"
)

var (
	c0, c1 *ClientImpl
)

func TestMain(m *testing.M) {
	c0 = &ClientImpl{}
	c0.buf = make(chan *common.Message, bufsize)
	c0.RPCServer = &RPCServer{c0.buf}
	serverAddr := "127.0.0.1:10000"
	c0.RunServer(serverAddr)

	c1 = &ClientImpl{}
	c1.ConnectPeers([]string{serverAddr})
	m.Run()
}

func TestRPCServer_Receive(t *testing.T) {
	msg := &common.Message{
		SenderID: 1,
		Body:     "test message1",
	}
	var reply string
	c1.peers[0].Call("RPCServer.Receive", msg, &reply)

	result := <-c0.buf
	if !reflect.DeepEqual(*result, *msg) {
		t.Errorf("\nexpected: %v\nactual: %v\n", *msg, *result)
	}
}

func TestClientImpl_Broadcast(t *testing.T) {
	msgs := []*common.Message{
		{
			SenderID: 1,
			Body:     "test message2",
		},
		{
			SenderID: 1,
			Body:     "test message3",
		},
	}
	for _, msg := range msgs {
		c1.Broadcast(msg)
	}

	for i := 0; i < len(msgs); i++ {
		result := c0.GetData()
		if !reflect.DeepEqual(*msgs[i], *result) {
			t.Errorf("\nexpected: %v\nactual: %v\n", *msgs[i], *result)
		}
	}
}
