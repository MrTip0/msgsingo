package main

import (
	"bufio"
	"fmt"
	"msgserver/mycrypto"
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

	key, err := mycrypto.GetKey(conn)
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
		name, err = mycrypto.GcmDecrypter(key, name)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err.Error())
			name = "Unknown name"
		}

		if *pass != "" {
			mycrypto.SendEncrypted(conn, "pass", key)
			insertedPass, err := mycrypto.ReadEncrypted(reader, key)
			if err != nil {
				return
			}
			if insertedPass != *pass {
				mycrypto.SendEncrypted(conn, "wrong", key)
				return
			} else {
				mycrypto.SendEncrypted(conn, "ok", key)
			}
		} else {
			mycrypto.SendEncrypted(conn, "nopass", key)
		}
	}

	updates <- update{rx: receiveChannel, action: Add, name: name}

	fmt.Printf("%s connected\n", addr)

	end := new(sync.Mutex)
	end.Lock()

	go receiveMsgs(conn, end, sendChannel, key, receiveChannel, name)
	go sendMsgs(conn, receiveChannel, key, terminate)

	end.Lock()
}

func receiveMsgs(conn net.Conn, end *sync.Mutex, sendChannel chan message, key []byte, toAvoid chan message, name string) {
	var data string
	var err error
	defer func() {
		end.Unlock()
	}()

	reader := bufio.NewReader(conn)

	message := message{"", name, toAvoid}

	for {
		data, err = mycrypto.ReadEncrypted(reader, key)
		if err != nil {
			return
		}

		fmt.Printf("INFO: received %s from %s\n", data, conn.RemoteAddr())
		message.message = data
		sendChannel <- message
		conn.SetDeadline(time.Now().Add(time.Minute * 5))

	}
}

func sendMsgs(conn net.Conn, rx chan message, key []byte, terminate chan bool) {
	for {
		select {
		case data := <-rx:
			if err := mycrypto.SendEncrypted(conn, fmt.Sprintf("%s - %s", data.author, data.message), key); err != nil {
				fmt.Printf("%s\n", err.Error())
				continue
			}

			conn.SetDeadline(time.Now().Add(time.Minute * 5))
		case <-terminate:
			return
		}
	}
}
