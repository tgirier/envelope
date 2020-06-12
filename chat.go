// Package chat implements a chat server & client
package chat

import (
	"log"
	"net"
	"sync"
)

// Server represents a chat server
type Server struct {
	sync.Mutex
	listener net.Listener
	running  bool
}

// Client represents a chat client
type Client struct {
	connection net.Conn
}

// Start returns a pointer to a running server
func Start(addr string) (*Server, error) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return &Server{}, err
	}

	s := &Server{
		listener: ln,
		running:  true,
	}

	go s.Run()
	return s, nil
}

// Run implements the logic handling connections
func (s *Server) Run() {
	conn, err := s.listener.Accept()
	if err != nil {
		log.Printf("connection failed: %v", err)
		return
	}
	if !s.Running() {
		conn.Close()
		return
	}
}

// Running indicates if the server can accept connections
func (s *Server) Running() bool {
	s.Lock()
	defer s.Unlock()
	return s.running
}

// Stop stops a running server
func (s *Server) Stop() {
	s.Lock()
	defer s.Unlock()
	s.running = false
	s.listener.Close()
}

// Connect returns a new client with a connection to the server
func Connect(addr string) (*Client, error) {
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
