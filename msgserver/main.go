package main

import (
	"bufio"
	"container/list"
	"fmt"
	"net"
	"time"
)

const (
	Aggiungi int = 0
	Rimuovi      = 1
)

type update struct {
	rx     chan string
	action int
}

func main() {
	ln, err := net.Listen("tcp", ":4242")
	if err != nil {
		fmt.Printf("An error occurred: %s\n", err.Error())
		return
	}

	txs := list.New()
	updates := list.New()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Fprintf(conn, "An error occurred: %s\n", err.Error())
			continue
		}

		fmt.Printf("%s has connected\n", conn.RemoteAddr())

		ntx := make(chan string, 10)
		updc := make(chan update)

		go handleConnection(conn, ntx, txs, updc)

		el := updates.Front()
		for el != nil {
			el.Value.(chan update) <- update{rx: ntx, action: Aggiungi}
			el = el.Next()
		}

		updates.PushBack(updc)

		txs.PushBack(ntx)
	}
}

func handleConnection(conn net.Conn, rx chan string, rxso *list.List, updates chan update) {
	conn.SetDeadline(time.Now().Add(time.Minute * 5))
	end1, end2 := make(chan bool, 1), make(chan bool, 1)

	others := list.New()

	{
		tmpl := rxso.Front()
		for tmpl != nil {
			if tmpl.Value != rx {
				others.PushBack(tmpl.Value)
			}
			tmpl = tmpl.Next()
		}
	}

	go receiveMsgs(conn, &others, end1)
	go sendMsgs(conn, rx, end2)

updcycle:
	for {
		select {
		case _ = <-end1:
			updates <- update{rx: rx, action: Rimuovi}
			break updcycle
		case _ = <-end2:
			updates <- update{rx: rx, action: Rimuovi}
			break updcycle
		default:
		}
		nel := <-updates
		if nel.action == Aggiungi {
			if nel.rx != rx {
				others.PushBack(nel.rx)
			}
		} else if nel.action == Rimuovi {
			el := others.Front()
			for el != nil && el.Value != nel.rx {
				el = el.Next()
			}
			if el != nil {
				others.Remove(el)
			}
		}
	}
}

func receiveMsgs(conn net.Conn, tx **list.List, end chan bool) {
	defer func() {
		conn.Close()
		end <- true
	}()

	reader := bufio.NewReader(conn)

	for {
		var data string
		data, err := reader.ReadString('\x00')
		if err != nil {
			return
		}
		if len(data) > 0 {
			conn.SetDeadline(time.Now().Add(time.Minute * 5))
			// Remove x00
			data = data[:len(data)-1]
			fmt.Printf("INFO: received %s from %s\n", data, conn.RemoteAddr())
			l := *tx
			el := l.Front()
			for el != nil {
				el.Value.(chan string) <- data
				el = el.Next()
			}
		}
	}
}

func sendMsgs(conn net.Conn, rx chan string, end chan bool) {
	defer func() {
		conn.Close()
		end <- true
	}()

	for {
		data := <-rx
		_, err := fmt.Fprintf(conn, "%s\x00", data)
		if err != nil {
			return
		}
		conn.SetDeadline(time.Now().Add(time.Minute * 5))
	}
}
