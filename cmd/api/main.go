package main

import (
	"log"

	"github.com/abc_valera/flugo/internal/server"
)

func main() {
	s, err := server.NewServer()
	if err != nil {
		log.Fatal(err)
	}

	s.Start()
}
