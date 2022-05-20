package main

import (
	"bufio"
	"container/list"
	"fmt"
	"net"
	"time"
)

func main() {
	ln, err := net.Listen("tcp", ":4242")
	if err != nil {
		fmt.Printf("An error occurred: %s\n", err.Error())
		return
	}

	rxs := list.New()
	txs := list.New()

	go sendMsgsBetweenProcesses(&txs, &rxs)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Fprintf(conn, "An error occurred: %s\n", err.Error())
			continue
		}

		fmt.Printf("%s has connected\n", conn.RemoteAddr())

		nrx := make(chan string, 10)
		ntx := make(chan string, 10)

		go handleConnection(conn, ntx, nrx)

		rxs.PushBack(nrx)
		txs.PushBack(ntx)
	}
}

func handleConnection(conn net.Conn, rx chan string, tx chan string) {
	conn.SetDeadline(time.Now().Add(time.Minute * 5))
	end1, end2 := make(chan bool, 1), make(chan bool, 1)

	go receiveMsgs(conn, tx, end1)
	go sendMsgs(conn, rx, end2)

	<-end1
	<-end2
}

func receiveMsgs(conn net.Conn, tx chan string, end chan bool) {
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
			tx <- data
		}
	}
}

func sendMsgs(conn net.Conn, rx chan string, end chan bool) {
	defer func() {
		conn.Close()
		end <- true
	}()

	for {
		select {
		case data := <-rx:
			_, err := fmt.Fprintf(conn, "%s\x00", data)
			if err != nil {
				return
			}
			conn.SetDeadline(time.Now().Add(time.Minute * 5))
		default:
		}
		time.Sleep(time.Millisecond * 10)
	}
}

func sendMsgsBetweenProcesses(txs **list.List, rxs **list.List) {
	for {
		listrx := *rxs
		listrxel := listrx.Front()
		for listrxel != nil {
			select {
			case data := <-listrxel.Value.(chan string):
				listtx := *txs
				listtxel := listtx.Front()
				for listtxel != nil {
					listtxel.Value.(chan string) <- data
					listtxel = listtxel.Next()
				}
			default:
			}
			listrxel = listrxel.Next()
		}
		time.Sleep(time.Millisecond * 10)
	}
}
