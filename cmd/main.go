package main

import (
	"fmt"
	"log"

	"github.com/tgirier/chat"
)

func main() {
	s := chat.NewServer()
	s.Port = 8080
	s.Host = "127.0.0.1"

	fmt.Println(s.ListenAddress())

	err := s.ListenAndServe()

	if err != nil {
		log.Fatalln("failed starting server")
	}
}
