package chat_test

import (
	"log"
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/phayes/freeport"
	"github.com/tgirier/chat"
)

func TestServerConn(t *testing.T) {
	t.Parallel()

	errChan := make(chan error)
	logger := newTestLogger(t)

	p, err := freeport.GetFreePort()
	if err != nil {
		t.Fatal("no port available")
	}
	addr := net.JoinHostPort("localhost", strconv.Itoa(p))
	go func() {
		errChan <- chat.ListenAndServeWithLogger(addr, logger)
	}()

	_, err = chat.ConnectClient(addr)
	if err != nil {
		t.Fatalf("client can't connect: %v", err)
	}
}

func TestUsername(t *testing.T) {
	t.Parallel()

	c := startServerAndClient(t)

	c.Read()

	username := "My Name"
	want := "My Name joined the chat\n"

	c.Send(username + "\n")
	got, err := c.Read()

	if err != nil {
		t.Fatalf("reading back the joined chat message failed: %v", err)
	}
	if got != want {
		t.Errorf("username sent: got %q, want %q", got, want)
	}
}

func TestSendMessageAndEcho(t *testing.T) {
	t.Parallel()

	c := startServerAndClient(t)

	c.Read()

	username := "My Name"
	msg := "Hello all"
	want := username + ": " + msg + "\n"

	c.Send(username + "\n")
	_, err := c.Read()

	c.Send(msg + "\n")
	got, err := c.Read()

	c.Close()

	// Enable server to log client closing
	time.Sleep(10 * time.Millisecond)

	if err != nil {
		t.Fatalf("reading back our own message failed:  %v", err)
	}
	if got != want {
		t.Errorf("sent message: got %q, want %q", got, want)
	}

}

func startServerAndClient(t *testing.T) *chat.Client {
	errChan := make(chan error)

	logger := newTestLogger(t)

	p, err := freeport.GetFreePort()
	if err != nil {
		t.Fatal("no port available")
	}
	addr := net.JoinHostPort("localhost", strconv.Itoa(p))

	go func() {
		errChan <- chat.ListenAndServeWithLogger(addr, logger)
	}()

	c, err := chat.ConnectClient(addr)
	if err != nil {
		t.Fatalf("client connection failed: %v", err)
	}
	return c
}

func newTestLogger(t *testing.T) *log.Logger {
	t.Helper()
	return log.New(testWriter{t}, t.Name()+" ", 0)
}

type testWriter struct {
	*testing.T
}

func (tw testWriter) Write(p []byte) (int, error) {
	tw.Helper()
	tw.Logf("%s", p)
	return len(p), nil
}
