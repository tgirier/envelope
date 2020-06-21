package chat_test

import (
	"testing"

	"github.com/tgirier/chat"
)

func TestServerConn(t *testing.T) {
	t.Parallel()

	s, err := chat.RandomPortServer()
	if err != nil {
		t.Fatal(err)
	}
	defer s.Stop()

	c, err := chat.ConnectClient(s.ListenAddress)
	if err != nil {
		t.Fatalf("connection failed: %v", err)
	}
	defer c.Close()
}

func TestServerClose(t *testing.T) {
	t.Parallel()

	s, err := chat.RandomPortServer()
	if err != nil {
		t.Fatal(err)
	}

	s.Stop()

	c, err := chat.ConnectClient(s.ListenAddress)
	if err == nil {
		t.Error("server still running")
		defer c.Close()
	}

}

func TestWelcomeMessage(t *testing.T) {
	t.Parallel()

	want := "Welcome to ChatRoom !"

	s, err := chat.RandomPortServer()
	if err != nil {
		t.Fatal(err)
	}
	defer s.Stop()

	c, err := chat.ConnectClient(s.ListenAddress)
	if err != nil {
		t.Fatalf("client connection failed: %v", err)
	}
	defer c.Close()

	got, err := c.Read()

	if err != nil {
		t.Fatalf("reading welcome message failed:  %v", err)
	}
	if got != want {
		t.Errorf("welcome message: got %q, want %q", got, want)
	}
}

// func TestAtoBMessage(t *testing.T) {
// 	t.Parallel()

// 	addr := ":8080"
// 	m := "Hello B !"
// 	want := "Hello B !"

// 	// Setting up connA and connB
// 	connA, err := net.Dial("tcp", addr)
// 	if err != nil {
// 		t.Fatalf("connection a failed: %v", err)
// 	}
// 	defer connA.Close()

// 	connB, err := net.Dial("tcp", addr)
// 	if err != nil {
// 		t.Fatalf("connection b failed: %v", err)
// 	}
// 	defer connB.Close()

// 	// Reading welcome message on connB and discarding it
// 	var b bytes.Buffer
// 	_, err = io.Copy(&b, connB)
// 	if err != nil {
// 		t.Fatalf("reading welcome message on connection b failed: %v", err)
// 	}

// 	b.Reset()

// 	// Sending Message on connA
// 	_, err = connA.Write([]byte(m))
// 	if err != nil {
// 		t.Fatalf("sending message on connection a failed: %v", err)
// 	}

// 	// Reading message on connB
// 	_, err = io.Copy(&b, connB)
// 	if err != nil {
// 		t.Fatalf("reading message on connection b failed: %v", err)
// 	}

// 	got := b.String()

// 	if got != want {
// 		t.Errorf("connection b received message: got %q, want %q", got, want)
// 	}

// }

// func TestBiDirectionnalMessages(t *testing.T) {
// 	t.Parallel()

// 	addr := ":8080"
// 	mA := "Hello B !"
// 	mB := "Hello A !"
// 	wantA := "Hello A !"
// 	wantB := "Hello B !"

// 	// Setting up connA and connB
// 	connA, err := net.Dial("tcp", addr)
// 	if err != nil {
// 		t.Fatalf("connection a failed: %v", err)
// 	}
// 	defer connA.Close()

// 	connB, err := net.Dial("tcp", addr)
// 	if err != nil {
// 		t.Fatalf("connection b failed: %v", err)
// 	}
// 	defer connB.Close()

// 	// Receiving welcome messages and discarding it
// 	var bA bytes.Buffer
// 	_, err = io.Copy(&bA, connA)
// 	if err != nil {
// 		t.Fatalf("reading welcome message on connection a failed: %v", err)
// 	}
// 	bA.Reset()

// 	var bB bytes.Buffer
// 	_, err = io.Copy(&bB, connB)
// 	if err != nil {
// 		t.Fatalf("reading welcome message on connection b failed: %v", err)
// 	}
// 	bB.Reset()

// 	// Sending messages on both connections
// 	_, err = connA.Write([]byte(mA))
// 	if err != nil {
// 		t.Fatalf("sending message on connection a failed: %v", err)
// 	}

// 	_, err = connB.Write([]byte(mB))
// 	if err != nil {
// 		t.Fatalf("sending message on connection b failed: %v", err)
// 	}

// 	// Receiving messages on both connections
// 	_, err = io.Copy(&bA, connA)
// 	if err != nil {
// 		t.Fatalf("receiving b message on connection a failed: %v", err)
// 	}

// 	_, err = io.Copy(&bB, connB)
// 	if err != nil {
// 		t.Fatalf("receiving a message on connection b failed: %v", err)
// 	}

// 	gotA := bA.String()
// 	gotB := bB.String()

// 	if gotA != wantA {
// 		t.Errorf("connection A message received: got %q, want %q", gotA, wantA)
// 	}
// 	if gotB != wantB {
// 		t.Errorf("connection B message received: got %q, want %q", gotB, wantB)
// 	}
// }
