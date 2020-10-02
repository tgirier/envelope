package envelope

import (
	"bufio"
	"fmt"
	"net"
)

// Client represents a chat client
type Client struct {
	connection net.Conn
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
