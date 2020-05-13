package chat_test

import (
	"bytes"
	"io"
	"net"
	"testing"
)

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
