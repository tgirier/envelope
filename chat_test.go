package chat_test

import (
	"log"
	"testing"
	"time"

	"github.com/tgirier/chat"
)

type myLogger struct {
}

func (l *myLogger) Println(v ...interface{}) {
	log.Printf("log: %s", v)
}
func TestServerConn(t *testing.T) {
	t.Parallel()

	s := chat.NewServer()
	s.Logger = &myLogger{} // look at logrus

	errChan := make(chan error)
	done := make(chan struct{})

	go func() {
		errChan <- s.ListenAndServe()
	}()
	defer s.Close()

	go func() {
		c, err := chat.ConnectClient(s.ListenAddress)
		if err != nil {
			errChan <- err
			return
		}
		defer c.Close()
		close(done)
	}()

	select {
	case err := <-errChan:
		t.Fatalf("connection failed: %v", err)
	case <-done:
		return
	}
}

func TestServerClose(t *testing.T) {
	t.Parallel()

	s := chat.NewServer()

	errChan := make(chan error)
	runningChan := make(chan struct{})

	go func() {
		errChan <- s.ListenAndServe()
	}()

	go func() {
		for !s.Running() {
			time.Sleep(10 * time.Millisecond)
		}
		close(runningChan)
	}()

	select {
	case err := <-errChan:
		t.Fatalf("starting server failed: %v", err)
	case <-runningChan:
		s.Close()
	}

	c, err := chat.ConnectClient(s.ListenAddress)
	if err == nil {
		t.Error("server still running")
		defer c.Close()
	}
}

func TestWelcomeMessage(t *testing.T) {
	t.Parallel()

	want := "Welcome to ChatRoom !\n"

	errChan := make(chan error)
	runningChan := make(chan struct{})

	s := chat.NewServer()

	go func() {
		errChan <- s.ListenAndServe()
	}()
	defer s.Close()

	go func() {
		for !s.Running() {
			time.Sleep(10 * time.Millisecond)
		}
		close(runningChan)
	}()
	// s, err := chat.StartServer()
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// defer s.Stop()

	select {
	case err := <-errChan:
		t.Fatalf("starting server failed: %v", err)
	case <-runningChan:
	}

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
	// t.Parallel()

	// want := "Welcome to ChatRoom !\n"

	// s, err := chat.StartServer()
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// defer s.Stop()

	// c, err := chat.ConnectClient(s.ListenAddress)
	// if err != nil {
	// 	t.Fatalf("client connection failed: %v", err)
	// }
	// defer c.Close()

	// got, err := c.Read()

	// if err != nil {
	// 	t.Fatalf("reading welcome message failed:  %v", err)
	// }
	// if got != want {
	// 	t.Errorf("welcome message: got %q, want %q", got, want)
	// }
}

// func TestSendMessageAndEcho(t *testing.T) {
// 	t.Parallel()

// 	s, c := startServerAndClient(t)
// 	defer s.Stop()
// 	defer c.Close()

// 	c.Read()

// 	want := "Hello all\n"

// 	c.Send(want)
// 	got, err := c.Read() // check for loop

// 	if err != nil {
// 		t.Fatalf("reading welcome message failed:  %v", err)
// 	}
// 	if got != want {
// 		t.Errorf("sent message: got %q, want %q", got, want)
// 	}

// }

// func TestMultipleAndEcho(t *testing.T) {
// 	t.Parallel()

// 	s, c := startServerAndClient(t)
// 	defer s.Stop()
// 	defer c.Close()

// 	c.Read()

// 	m1 := "Hello all\n"
// 	want := "Second message\n"

// 	c.Send(m1)
// 	// fmt.Println("client: message 1 sent")
// 	c.Read()
// 	// fmt.Printf("client: message 1 received %s", m) //Check for debug method
// 	c.Send(want)
// 	// fmt.Println("client: message 2 sent")
// 	got, err := c.Read()

// 	if err != nil {
// 		t.Fatalf("reading welcome message failed:  %v", err)
// 	}
// 	if got != want {
// 		t.Errorf("sent message: got %q, want %q", got, want)
// 	}

// }

// func startServerAndClient(t *testing.T) (*chat.Server, *chat.Client) {
// 	s, err := chat.StartServer()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	c, err := chat.ConnectClient(s.ListenAddress)
// 	if err != nil {
// 		t.Fatalf("client connection failed: %v", err)
// 	}
// 	return s, c
// }

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
