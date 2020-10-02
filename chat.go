// Package chat implements a chat server & client
package chat

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

// server represents a chat server
type server struct {
	logger         logrus.StdLogger
	welcomeMessage string
	listener       net.Listener
	clients        map[*net.Conn]string
	register       chan *connection
	unregister     chan *net.Conn
	broadcast      chan string
}

// Connection represents a user connection
type connection struct {
	conn     *net.Conn
	username string
}

// Run implements the logic handling connections
func (s *server) run() {

	go s.listen()

	for {
		select {
		case c := <-s.register:
			s.clients[c.conn] = c.username
			s.logger.Println("client connection registered")
			msg := fmt.Sprintf("%s joined envelope\n", c.username)
			s.broadcast <- msg
			go s.handle(c)
		case msg := <-s.broadcast:
			for conn := range s.clients {
				_, err := fmt.Fprint(*conn, msg)
				if err != nil {
					s.logger.Println(fmt.Sprintf("sending message failed: %v", err))
				}
			}
		case conn := <-s.unregister:
			delete(s.clients, conn)
			s.logger.Println("client connection unregistered")
		}
	}
}

func (s *server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			s.logger.Println(fmt.Sprintf("connection failed: %v", err))
			continue
		}

		_, err = fmt.Fprintf(conn, s.welcomeMessage+"\n")
		if err != nil {
			s.logger.Println(fmt.Sprintf("sending message failed: %v", err))
			continue
		}

		r := bufio.NewReader(conn)

		username, err := r.ReadString('\n')
		if err != nil {
			s.logger.Println(fmt.Sprintf("reading username failed: %v", err))
			continue
		}

		username = strings.TrimSuffix(username, "\n")

		c := &connection{
			conn:     &conn,
			username: username,
		}

		s.register <- c
	}
}

func (s *server) handle(c *connection) {
	r := bufio.NewReader(*c.conn)

	for {
		msg, err := r.ReadString('\n')
		if err == io.EOF {
			s.logger.Println("client connection closed")
			s.unregister <- c.conn
			break
		}
		if err != nil {
			s.logger.Println(fmt.Sprintf("receiving message failed: %v", err))
		}

		msg = c.username + ": " + msg

		s.broadcast <- msg
	}
}

// ListenAndServe blocks while the server is running
func ListenAndServe(addr string) (err error) {
	return ListenAndServeWithLogger(addr, log.New(os.Stderr, "", log.LstdFlags))
}

// ListenAndServeWithLogger blocks while the server is running.
// Enable logger customization
func ListenAndServeWithLogger(addr string, logger logrus.StdLogger) (err error) {
	s := &server{
		welcomeMessage: "Welcome to envelope! Please enter your username:",
	}

	s.logger = logger
	s.register = make(chan *connection, 1)
	s.unregister = make(chan *net.Conn, 1)
	s.clients = make(map[*net.Conn]string)
	s.broadcast = make(chan string, 10)

	s.listener, err = net.Listen("tcp", addr)

	if err != nil {
		return err
	}

	s.logger.Println(fmt.Sprintf("Listening on %v", addr))

	s.run()

	return nil
}
