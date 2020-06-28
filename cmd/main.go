package main

import (
	"fmt"
	"log"

	"github.com/tgirier/chat"
)

func main() {
	s, err := chat.StartServer(chat.WithPort(8080))
	if err != nil {
		log.Fatalln("failed starting server")
	}
	fmt.Println(s.ListenAddress)

	s.ListenAndServe()
}
