package main

import (
	"log"

	"github.com/tgirier/envelope"
)

func main() {
	log.Fatal(envelope.ListenAndServe("localhost:8080"))
}
