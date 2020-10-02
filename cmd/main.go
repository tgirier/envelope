package main

import (
	"log"

	"github.com/tgirier/chat"
)

func main() {
	log.Fatal(chat.ListenAndServe("localhost:8080"))
}
