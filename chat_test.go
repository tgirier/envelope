package chat_test

import (
	"log"
	"net"
	"strconv"
	"testing"

	"github.com/phayes/freeport"
	"github.com/tgirier/chat"
)

type myLogger struct {
}

func (l *myLogger) Println(v ...interface{}) {
	log.Printf("log: %s", v)
}
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

// func TestServerClose(t *testing.T) {
// 	t.Parallel()

// 	errChan := make(chan error)

// 	p, err := freeport.GetFreePort()
// 	if err != nil {
// 		t.Fatal("no port available")
// 	}
// 	addr := net.JoinHostPort("localhost", strconv.Itoa(p))

// 	go func() {
// 		errChan <- chat.ListenAndServe(addr)
// 	}()

// 	c, err := chat.ConnectClient(addr)
// 	if err == nil {
// 		t.Error("server still running")
// 		defer c.Close()
// 	}
// }

// func TestPortSwitching(t *testing.T) {
// 	t.Parallel()

// 	s1 := chat.NewServer()
// 	s1.Port = 8080
// 	s2 := chat.NewServer()
// 	s2.Port = 8080

// 	errChan := make(chan error)
// 	runningChan := make(chan struct{})

// 	go func() {
// 		errChan <- s1.ListenAndServe()
// 	}()
// 	defer s1.Close()

// 	go func() {
// 		errChan <- s2.ListenAndServe()
// 	}()
// 	defer s2.Close()

// 	go func() {
// 		for !s1.Running() && !s2.Running() {
// 			time.Sleep(10 * time.Millisecond)
// 		}
// 		close(runningChan)
// 	}()

// 	select {
// 	case err := <-errChan:
// 		t.Fatalf("failed starting server: %v", err)
// 	case <-runningChan:
// 	}

// 	if s1.Port == s2.Port {
// 		t.Errorf("switching port failed: s1 port %d, s2 port %d", s1.Port, s2.Port)
// 	}

// }

// func TestWelcomeMessage(t *testing.T) {
// 	t.Parallel()

// 	want := "Welcome to Thibaut's chat !\n"

// 	errChan := make(chan error)

// 	s := chat.NewServer()
// 	s.WelcomeMessage = want

// 	p, err := freeport.GetFreePort()
// 	if err != nil {
// 		t.Fatal("no port available")
// 	}
// 	addr := net.JoinHostPort("localhost", strconv.Itoa(p))

// 	go func() {
// 		errChan <- chat.ListenAndServe(addr)
// 	}()
// 	defer s.Close()

// 	<-s.Ready

// 	c, err := chat.ConnectClient(addr)
// 	if err != nil {
// 		t.Fatalf("client connection failed: %v", err)
// 	}
// 	defer c.Close()

// 	got, err := c.Read()

// 	if err != nil {
// 		t.Fatalf("reading back our own message failed:  %v", err)
// 	}
// 	if got != want {
// 		t.Errorf("welcome message: got %q, want %q", got, want)
// 	}
// }

func TestSendMessageAndEcho(t *testing.T) {
	t.Parallel()

	c := startServerAndClient(t)
	defer c.Close()

	c.Read()

	want := "Hello all\n"

	c.Send(want)
	got, err := c.Read() // check for loop

	if err != nil {
		t.Fatalf("reading back our own message failed:  %v", err)
	}
	if got != want {
		t.Errorf("sent message: got %q, want %q", got, want)
	}

}

func TestMultipleAndEcho(t *testing.T) {
	t.Parallel()

	c := startServerAndClient(t)
	defer c.Close()

	c.Read()

	m1 := "Hello all\n"
	want := "Second message\n"

	c.Send(m1)
	// fmt.Println("client: message 1 sent")
	c.Read()
	// fmt.Printf("client: message 1 received %s", m) //Check for debug method
	c.Send(want)
	// fmt.Println("client: message 2 sent")
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
