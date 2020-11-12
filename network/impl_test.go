package network

import (
	"os"
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
	peer := &Peer{ID: 0, Address: serverAddr, Client: nil}
	c1.ConnectPeers([]*Peer{peer})
	m.Run()
	os.Exit(0)
}

func TestRPCServer_Receive(t *testing.T) {
	msg := &common.Message{
		SenderID: 1,
		Body:     "test message1",
	}
	var reply string
	err := c1.peers[0].Client.Call("RPCServer.Receive", msg, &reply)
	if err != nil {
		t.Error(err)
	}

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

func TestClientImpl_GetPeers(t *testing.T) {
	c := &ClientImpl{}
	c.peers = []*Peer{
		{ID: 0, Address: "127.0.0.1:10000", Client: nil},
		{ID: 1, Address: "127.0.0.1:10000", Client: nil},
		{ID: 2, Address: "127.0.0.1:10000", Client: nil},
	}
	expected := []int{0, 1, 2}
	result := c.GetPeers()
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("\nexpected: %v\nactual: %v\n", expected, result)
	}
}
