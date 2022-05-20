package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

func main() {
	var serverhost string

	fmt.Printf("Insert the server ip:\n> ")
	fmt.Scanf("%s", &serverhost)

	port := "4242"

	if strings.Contains(serverhost, ":") {
		expl := strings.Split(serverhost, ":")
		serverhost, port = expl[0], expl[1]
	}
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", serverhost, port))

	if err != nil {
		fmt.Printf("An error occurred: %s\n", err.Error())
		return
	}
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(time.Minute * 5))

	fmt.Printf("Connecting...\n")
	password := getPass(conn)
	fmt.Printf("Connected\n")

	end := make(chan bool, 1)

	go sender(conn, password)
	go receiver(conn, end, password)

	<-end
}

func sender(conn net.Conn, pass []byte) {
	var msg string
	var err error

	reader := bufio.NewReader(os.Stdin)

	for {
		msg, err = reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error while reading input: %s\n", err.Error())
			continue
		}

		data, err := gcmEncrypter(pass, msg[:len(msg)-1])
		if err != nil {
			fmt.Printf("Error while encripting input: %s\n", err.Error())
			continue
		}

		_, err = fmt.Fprintf(conn, "%s\x00", data)
		if err != nil {
			fmt.Printf("Error while sending: %s\n", err.Error())
			return
		}
		conn.SetDeadline(time.Now().Add(time.Minute * 5))
	}
}

func receiver(conn net.Conn, end chan bool, pass []byte) {
	defer func() {
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
			data = data[:len(data)-1]
			data, err = gcmDecrypter(pass, data)
			if err != nil {
				fmt.Printf("error: %s\n", err.Error())
				continue
			}
			fmt.Printf("received: %s\n", data)
		}
	}
}
