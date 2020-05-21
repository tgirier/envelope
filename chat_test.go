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
	_, err = io.Copy(&b, conn)

	got := b.String()

	if err != nil {
		t.Fatalf("reading welcome message failed:  %v", err)
	}
	if got != want {
		t.Errorf("welcome message: got %q, want %q", got, want)
	}
}

func TestAtoBMessage(t *testing.T) {
	t.Parallel()

	addr := ":8080"
	m := "Hello B !"
	want := "Hello B !"

	// Setting up connA and connB
	connA, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatalf("connection a failed: %v", err)
	}
	defer connA.Close()

	connB, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatalf("connection b failed: %v", err)
	}
	defer connB.Close()

	// Reading welcome message on connB and discarding it
	var b bytes.Buffer
	_, err = io.Copy(&b, connB)
	if err != nil {
		t.Fatalf("reading welcome message on connection b failed: %v", err)
	}

	b.Reset()

	// Sending Message on connA
	_, err = connA.Write([]byte(m))
	if err != nil {
		t.Fatalf("sending message on connection a failed: %v", err)
	}

	// Reading message on connB
	_, err = io.Copy(&b, connB)
	if err != nil {
		t.Fatalf("reading message on connection b failed: %v", err)
	}

	got := b.String()

	if got != want {
		t.Errorf("connection b received message: got %s, want %s", got, want)
	}

}
