package chat_test

import (
	"net"
	"strconv"
	"testing"

	"github.com/phayes/freeport"
	"github.com/tgirier/chat"
)

func TestServerConn(t *testing.T) {
	t.Parallel()

	errChan := make(chan error)

	p, err := freeport.GetFreePort()
	if err != nil {
		t.Fatal("no port available")
	}
	addr := net.JoinHostPort("localhost", strconv.Itoa(p))
	go func() {
		errChan <- chat.ListenAndServe(addr)
	}()

	// runtime.Gosched()

	c, err := chat.ConnectClient(addr)
	if err != nil {
		t.Fatalf("client can't connect: %v", err)
	}
	defer c.Close()
}

func TestSendMessageAndEcho(t *testing.T) {
	t.Parallel()

	c := startServerAndClient(t)
	defer c.Close()

	c.Read()

	username := "My Name"
	msg := "Hello all"
	want := username + ": " + msg + "\n"

	c.Send(username + "\n")
	_, err := c.Read()

	c.Send(msg + "\n")
	got, err := c.Read()

	if err != nil {
		t.Fatalf("reading back our own message failed:  %v", err)
	}
	if got != want {
		t.Errorf("sent message: got %q, want %q", got, want)
	}

}

func startServerAndClient(t *testing.T) *chat.Client {
	errChan := make(chan error)

	p, err := freeport.GetFreePort()
	if err != nil {
		t.Fatal("no port available")
	}
	addr := net.JoinHostPort("localhost", strconv.Itoa(p))

	go func() {
		errChan <- chat.ListenAndServe(addr)
	}()

	c, err := chat.ConnectClient(addr)
	if err != nil {
		t.Fatalf("client connection failed: %v", err)
	}
	return c
}
