package chat

import (
	"log"
	"net"
	"sync"
)

type Server struct {
	sync.Mutex
	listener net.Listener
	running  bool
}

type Client struct {
	connection net.Conn
}

func Start(addr string) (Server, error) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return Server{}, err
	}

	s := Server{
		listener: ln,
		running:  true,
	}

	go s.Run()
	return s, nil
}

func (s *Server) Run() {
	conn, err := s.listener.Accept()
	if !s.Running() {
		conn.Close()
		return
	}
	if err != nil {
		log.Printf("connection failed: %v", err)
	}
}

func (s *Server) Running() bool {
	s.Lock()
	defer s.Unlock()
	return s.running
}

func (s *Server) Stop() {
	s.Lock()
	defer s.Unlock()
	s.running = false
	s.listener.Close()
}

func Connect(addr string) (Client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return Client{}, err
	}
	c := Client{
		connection: conn,
	}
	return c, nil
}

func (c *Client) Close() {

}
