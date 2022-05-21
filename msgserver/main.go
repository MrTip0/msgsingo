package main

import (
	"fmt"
	"net"
)

const (
	Add    int = 0
	Remove     = 1
)

type update struct {
	action int
	rx     chan message
	name   string
}

type message struct {
	message, author string
	rx              chan message
}

func main() {
	ln, err := net.Listen("tcp", ":4242")
	if err != nil {
		fmt.Printf("An error occurred: %s\n", err.Error())
		return
	}

	updChan := make(chan update, 10)
	sender := make(chan message, 100)

	go connectionController(updChan, sender)

	fmt.Printf("Server Ready\n")

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Fprintf(conn, "An error occurred: %s\n", err.Error())
			continue
		}

		fmt.Printf("%s is connecting\n", conn.RemoteAddr())

		newSender := make(chan message, 10)

		go handleConnection(conn, newSender, sender, updChan)
	}
}
