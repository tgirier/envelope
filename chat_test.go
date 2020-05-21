package chat_test

import (
	"bytes"
	"io"
	"log"
	"net"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	serverErrors := make(chan error, 1)

	go func() {
		addr := ":8080"
		serverErrors <- chat.Listen("tcp", addr)
	}()

	select {
	case err := <-serverErrors:
		log.Fatalf("starting server failed: %v", err)
	default:
		os.Exit(m.Run())
	}

}

func TestServerConn(t *testing.T) {
	t.Parallel()

	addr := ":8080"

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatalf("connection failed: %v", err)
	}
	conn.Close()
}

func TestWelcomeMessage(t *testing.T) {
	t.Parallel()

	addr := ":8080"
	want := "Welcome to ChatRoom !"

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatalf("connection failed: %v", err)
	}
	defer conn.Close()

	var b bytes.Buffer
	io.Copy(&b, conn)

	got := b.String()

	if err != nil {
		t.Errorf("reading welcome message failed:  %v", err)
	}
	if got != want {
		t.Errorf("welcome message: got %q, want %q", got, want)
	}
}
