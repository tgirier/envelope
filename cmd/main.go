package main

import (
	"fmt"
	"log"

	"github.com/tgirier/chat"
)

func main() {
	s := chat.NewServer()
	s.Port = 8080

	fmt.Println(s.ListenAddress())

	log.Fatal(s.ListenAndServe())
}
