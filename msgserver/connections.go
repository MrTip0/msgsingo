package main

import (
	"bufio"
	"container/list"
	"fmt"
	"net"
	"time"
)

func handleConnection(conn net.Conn, rx chan string, rxso *list.List, updates chan update, dies chan chan string) {
	addr := conn.RemoteAddr()
	defer func() {
		dies <- rx
		fmt.Printf("%s has disconnected\n", addr)
		conn.Close()
	}()

	conn.SetDeadline(time.Now().Add(time.Minute * 5))
	end := make(chan bool, 1)

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

	password, err := getPass(conn)
	if err != nil {
		return
	}
	fmt.Printf("%s connected\n", addr)

	go receiveMsgs(conn, &others, end, password)
	go sendMsgs(conn, rx, password)

	for {
		select {
		case _ = <-end:
			return
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

func receiveMsgs(conn net.Conn, tx **list.List, end chan bool, pass []byte) {
	defer func() {
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
			data, err = gcmDecrypter(pass, data)
			if err != nil {
				fmt.Printf("ERROR: %s\n", err.Error())
				continue
			}

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

func sendMsgs(conn net.Conn, rx chan string, pass []byte) {
	for {
		data := <-rx
		data, err := gcmEncrypter(pass, data)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err.Error())
		}

		_, err = fmt.Fprintf(conn, "%s\x00", data)
		if err != nil {
			return
		}
		conn.SetDeadline(time.Now().Add(time.Minute * 5))
	}
}
