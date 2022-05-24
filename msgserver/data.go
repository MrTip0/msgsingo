package main

import "flag"

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

// flags
var (
	pass *string
)

func initFlags() {
	pass = flag.String("password", "", "The password to allows user to connect")

	flag.Parse()
}
