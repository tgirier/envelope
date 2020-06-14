// Package chat implements a chat server & client
package chat

import (
	"fmt"
	"math/rand"
	"net"
	"sync"
	"time"
)

// Server represents a chat server
type Server struct {
	mutex         sync.Mutex
	listener      net.Listener
	running       bool
	Logger        Logger
	ListenAddress string
}

// Logger enables a customization of the log function
type Logger interface {
	Log(s string)
}

// StandardLogger defines a standard logger that implements the Logger interface
type StandardLogger struct {
	timeFormat string
}

// Client represents a chat client
type Client struct {
	connection net.Conn
}

// StartServer returns a pointer to a running server
func StartServer(addr string) (*Server, error) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return &Server{}, err
	}

	logger := NewStandardLogger(time.RFC3339)

	s := &Server{
		listener:      ln,
		running:       true,
		Logger:        logger,
		ListenAddress: addr,
	}

	go s.Run()
	return s, nil
}

// Run implements the logic handling connections
func (s *Server) Run() {
	conn, err := s.listener.Accept()
	if err != nil {
		s.Logger.Log(fmt.Sprintf("connection failed: %v\n", err))
		return
	}
	if !s.Running() {
		conn.Close()
		return
	}
}

// Running indicates if the server can accept connections
func (s *Server) Running() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.running
}

// Stop stops a running server
func (s *Server) Stop() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.running = false
	s.listener.Close()
}

// RandomPortServer returns a server listening on a random port
func RandomPortServer() (*Server, error) {
	rand.Seed(time.Now().UnixNano())

	p := 8080 + rand.Intn(10)
	addr := fmt.Sprintf("localhost:%d", p)

	return StartServer(addr)
}

// ConnectClient returns a new client with a connection to the server
func ConnectClient(addr string) (*Client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return &Client{}, err
	}
	c := &Client{
		connection: conn,
	}
	return c, nil
}

// Close close the connection between the client and the server
func (c *Client) Close() {
	c.connection.Close()
}

// Read reads message received by the client
func (c *Client) Read() (string, error) {
	return "", nil
}

// Log prints a standard log message
func (l *StandardLogger) Log(s string) {
	fmt.Printf("%s %s", time.Now().Format(l.timeFormat), s)
}

// NewStandardLogger returns a standard logger for the server
func NewStandardLogger(timeFormat string) *StandardLogger {
	return &StandardLogger{
		timeFormat: timeFormat,
	}
}
