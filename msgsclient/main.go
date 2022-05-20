package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:4242")
	if err != nil {
		fmt.Printf("An error occurred: %s\n", err.Error())
		return
	}
	conn.SetDeadline(time.Now().Add(time.Minute * 5))

	end1, end2 := make(chan bool, 1), make(chan bool, 1)

	go sender(conn, end1)
	go receiver(conn, end2)

	<-end1
	<-end2
}

func sender(conn net.Conn, end chan bool) {
	defer func() {
		conn.Close()
		end <- true
	}()

	var msg string
	var err error

	reader := bufio.NewReader(os.Stdin)

	for {
		msg, err = reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error while reading input: %s\n", err.Error())
			continue
		}
		_, err = fmt.Fprintf(conn, "%s\x00", msg[:len(msg)-1])
		if err != nil {
			fmt.Printf("Error while sending: %s\n", err.Error())
			return
		}
		conn.SetDeadline(time.Now().Add(time.Minute * 5))
	}
}

func receiver(conn net.Conn, end chan bool) {
	defer func() {
		conn.Close()
		end <- true
	}()

	reader := bufio.NewReader(conn)
	var data string
	var err error

	for {
		if data, err = reader.ReadString('\x00'); err != nil {
			return
		}
		conn.SetDeadline(time.Now().Add(time.Minute * 5))

		if len(data) > 0 && data != "\x00" {
			fmt.Printf("received: %s\n", data)
		}
	}
}
