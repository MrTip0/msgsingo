package main

import (
	"container/list"
	"fmt"
	"net"
)

const (
	Aggiungi int = 0
	Rimuovi      = 1
)

type update struct {
	action int
	rx     chan string
}

func main() {
	ln, err := net.Listen("tcp", ":4242")
	if err != nil {
		fmt.Printf("An error occurred: %s\n", err.Error())
		return
	}

	dies := make(chan chan string, 10)

	txs := list.New()
	updates := list.New()

	go func() {
		for {
			died := <-dies
			el := txs.Front()
			elup := updates.Front()

			for el.Value != died {
				el = el.Next()
				elup = elup.Next()
			}
			txs.Remove(el)
			updates.Remove(elup)

			el = updates.Front()
			for el != nil {
				el.Value.(chan update) <- update{Rimuovi, died}
			}
		}
	}()

	fmt.Printf("Server Ready\n")

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Fprintf(conn, "An error occurred: %s\n", err.Error())
			continue
		}

		fmt.Printf("%s is connecting\n", conn.RemoteAddr())

		ntx := make(chan string, 10)
		updc := make(chan update)

		go handleConnection(conn, ntx, txs, updc, dies)

		el := updates.Front()
		for el != nil {
			select {
			case el.Value.(chan update) <- update{rx: ntx, action: Aggiungi}:
				el = el.Next()
			default:
				prev := el
				el = el.Next()
				updates.Remove(prev)
			}
		}

		updates.PushBack(updc)

		txs.PushBack(ntx)
	}
}
