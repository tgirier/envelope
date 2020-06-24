// Package chat implements a chat server & client
package chat

import (
	"bufio"
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
	Logger        logger // with standard logger can be extended with logrus
	ListenAddress string
}

// Logger enables a customization of the log function
type logger interface {
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
		conn.Close() // Not sure if it is still useful as listener.close closes all connections
		return
	}
	_, err = conn.Write([]byte("Welcome to ChatRoom !\n"))
	if err != nil {
		s.Logger.Log(fmt.Sprintf("sending message failed: %v\n", err))
	}

	r := bufio.NewReader(conn)

	for s.Running() {
		m, err := r.ReadString('\n')
		if err != nil {
			s.Logger.Log(fmt.Sprintf("receiving message failed: %v\n", err))
		}
		// fmt.Printf("server: message received %q", m)
		_, err = fmt.Fprintf(conn, m)
		if err != nil {
			s.Logger.Log(fmt.Sprintf("sending message failed: %v\n", err))
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

// Stop stops a running server
func (s *Server) Stop() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.running = false
	s.listener.Close()
}

// ListenAndServe blocks while the server is running
func (s *Server) ListenAndServe() {
	for s.Running() {
		time.Sleep(5 * time.Second)
	}
}

// RandomPortServer returns a server listening on a random port
func RandomPortServer() (*Server, error) {
	rand.Seed(time.Now().UnixNano())

	p := 8080 + rand.Intn(20) // Add used port detection
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
	// Check for bufio.scanner
	r := bufio.NewReader(c.connection)

	m, err := r.ReadString('\n')
	if err != nil {
		return "", err
	}
	return m, nil
}

// Send sends message to the server
func (c *Client) Send(m string) error {
	_, err := fmt.Fprint(c.connection, m)
	return err
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
