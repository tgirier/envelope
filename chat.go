// Package chat implements a chat server & client
package chat

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"sync"
	"time"
)

// Server represents a chat server
type Server struct {
	Logger         Logger // with standard logger can be extended with logrus
	Host           string
	Port           int
	WelcomeMessage string
	mutex          sync.Mutex
	listener       net.Listener
	running        bool
}

// Logger enables a customization of the log function
type Logger interface {
	Println(v ...interface{})
}

// NewServer returns a pointer to a new server on localhost and random port
func NewServer() *Server {
	rand.Seed(time.Now().UnixNano())
	p := 49152 + rand.Intn(16383) // Add used port detection

	logger := log.New(os.Stderr, "", log.LstdFlags)

	s := &Server{
		Host:           "localhost",
		Port:           p,
		Logger:         logger,
		WelcomeMessage: "Welcome to Chat Room!",
		running:        false,
	}

	return s
}

// Run implements the logic handling connections
func (s *Server) Run() {
	conn, err := s.listener.Accept()
	if err != nil {
		s.Logger.Println(fmt.Sprintf("connection failed: %v", err))
		return
	}
	if !s.Running() {
		conn.Close() // Not sure if it is still useful as listener.close closes all connections
		return
	}
	_, err = conn.Write([]byte(s.WelcomeMessage + "\n"))
	if err != nil {
		s.Logger.Println(fmt.Sprintf("sending message failed: %v", err))
	}

	r := bufio.NewReader(conn)

	for s.Running() {
		m, err := r.ReadString('\n')
		if err == io.EOF {
			s.Logger.Println(fmt.Sprintf("client connection closed: %v", err))
			conn.Close()
			break
		}
		if err != nil {
			s.Logger.Println(fmt.Sprintf("receiving message failed: %v", err))
		}
		// fmt.Printf("server: message received %q", m)
		_, err = fmt.Fprintf(conn, m)
		if err != nil {
			s.Logger.Println(fmt.Sprintf("sending message failed: %v", err))
		}
		// fmt.Printf("server: message sent %q", m)
	}
}

// Running indicates if the server can accept connections
func (s *Server) Running() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.running
}

// Close closes all connection to the server
func (s *Server) Close() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.running = false
	s.listener.Close()
}

// ListenAndServe blocks while the server is running
func (s *Server) ListenAndServe() error {
	retry := 0
	ln, err := net.Listen("tcp", s.ListenAddress())
	for err != nil && retry < 2 {
		s.Logger.Println(fmt.Sprintf("port not available: %v", err))

		rand.Seed(time.Now().UnixNano())
		p := 49152 + rand.Intn(16383)

		s.mutex.Lock()
		s.Port = p
		s.mutex.Unlock()

		s.Logger.Println(fmt.Sprintf("Switching to port: %d", p))
		ln, err = net.Listen("tcp", s.ListenAddress())
		retry++
	}

	if err != nil {
		return err
	}

	s.mutex.Lock()
	s.listener = ln
	s.running = true
	s.mutex.Unlock()

	s.Run()

	return nil
}

// ListenAddress returns the address on which the server is listening
func (s *Server) ListenAddress() string {
	return fmt.Sprintf(s.Host+":%d", s.Port)
}
