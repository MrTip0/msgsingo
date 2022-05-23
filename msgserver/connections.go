package main

import (
	"bufio"
	"fmt"
	"net"
	"sync"
	"time"
)

// conn -> Connection;
// receiveChannel -> the channel where the function takes messages to be send;
// sendChannel -> the channel where the function send the received messages;
// updates -> the channel where the function sends when it is end and when the connection is created;
func handleConnection(conn net.Conn, receiveChannel, sendChannel chan message, updates chan update) {
	var name string = "unknown name"
	addr := conn.RemoteAddr()
	terminate := make(chan bool, 1)
	defer func() {
		fmt.Printf("%s has disconnected\n", addr)
		updates <- update{rx: receiveChannel, action: Remove, name: name}
		terminate <- true
		conn.Close()
	}()

	conn.SetDeadline(time.Now().Add(time.Minute * 5))

	end := new(sync.Mutex)
	end.Lock()

	password, err := getPass(conn)
	if err != nil {
		return
	}

	{
		reader := bufio.NewReader(conn)
		name, err = reader.ReadString('\x00')
		if err != nil {
			return
		}
		name = name[:len(name)-1]
		name, err = gcmDecrypter(password, name)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err.Error())
			name = "Unknown name"
		}
	}

	updates <- update{rx: receiveChannel, action: Add, name: name}

	fmt.Printf("%s connected\n", addr)

	go receiveMsgs(conn, end, sendChannel, password, receiveChannel, name)
	go sendMsgs(conn, receiveChannel, password, terminate)

	end.Lock()
}

func receiveMsgs(conn net.Conn, end *sync.Mutex, sendChannel chan message, pass []byte, toAvoid chan message, name string) {
	var data string
	var err error
	defer func() {
		end.Unlock()
	}()

	reader := bufio.NewReader(conn)

	message := message{"", name, toAvoid}

	for {
		data, err = reader.ReadString('\x00')
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
			message.message = data
			sendChannel <- message
		}
	}
}

func sendMsgs(conn net.Conn, rx chan message, pass []byte, terminate chan bool) {
	for {
		select {
		case data := <-rx:
			payload, err := gcmEncrypter(pass, fmt.Sprintf("%s - %s", data.author, data.message))
			if err != nil {
				fmt.Printf("ERROR: %s\n", err.Error())
			}

			_, err = fmt.Fprintf(conn, "%s\x00", payload)
			if err != nil {
				return
			}
			conn.SetDeadline(time.Now().Add(time.Minute * 5))
		case <-terminate:
			return
		}
	}
}
