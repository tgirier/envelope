package chat_test

import (
	"net"
	"testing"
)

func TestServerConn(t *testing.T) {
	addr := ":8080"

	conn, err := net.Dial("tcp", addr)
	defer conn.Close()

	if err != nil {
		t.Errorf("connection failed: %v", err)
	}
}

func TestWelcomeMessage(t *testing.T) {
	addr := ":8080"
	want := "Welcome to ChatRoom !"

	conn, _ := net.Dial("tcp", addr)
	defer conn.Close()

	buf := make([]byte, 512)
	_, err := conn.Read(buf)

	got := string(buf)

	if err != nil {
		t.Errorf("reading welcome message failed:  %v", err)
	}
	if got != want {
		t.Errorf("welcome message: got %q, want %q", got, want)
	}
}
