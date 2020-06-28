package main

import (
	"fmt"
	"log"

	"github.com/tgirier/chat"
)

func main() {
	s, err := chat.StartServer()
	if err != nil {
		log.Fatalln("failed starting server")
	}
	fmt.Println(s.ListenAddress)

	s.ListenAndServe()
}
