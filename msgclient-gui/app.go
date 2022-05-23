package main

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx     context.Context
	conn    net.Conn
	connKey []byte
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) shutdown(ctx context.Context) {
	ctx.Done()
	if a.conn != nil {
		a.conn.Close()
	}
}

func (a *App) CreateConnection(host, user string) {
	port := "4242"

	if strings.Contains(host, ":") {
		expl := strings.Split(host, ":")
		host, port = expl[0], expl[1]
	}
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", host, port))

	if err != nil {
		fmt.Printf("An error occurred: %s\n", err.Error())
		return
	}

	conn.SetDeadline(time.Now().Add(time.Minute * 5))

	password := getPass(conn)

	{
		data, err := gcmEncrypter(password, user)
		if err != nil {
			fmt.Printf("Error while encripting username: %s\n", err.Error())
		}
		fmt.Fprintf(conn, "%s\x00", data)
	}
	a.conn = conn
	a.connKey = password

	go receiver(conn, password, a.ctx)
	runtime.EventsEmit(a.ctx, "connectionOpened")
	fmt.Println("Connected")
}

func (a *App) SendMessage(msg string) {
	var err error

	data, err := gcmEncrypter(a.connKey, msg)
	if err != nil {
		fmt.Printf("Error while encripting input: %s\n", err.Error())
	}

	_, err = fmt.Fprintf(a.conn, "%s\x00", data)
	if err != nil {
		fmt.Printf("Error while sending: %s\n", err.Error())
		return
	}
}

func receiver(conn net.Conn, pass []byte, ctx context.Context) {
	defer func() {
		runtime.EventsEmit(ctx, "connectionClosed")
	}()

	reader := bufio.NewReader(conn)

	var data string
	var err error

	for {
		if data, err = reader.ReadString('\x00'); err != nil {
			fmt.Println("The server has disconnected")
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
			runtime.EventsEmit(ctx, "receivedMessage", data)
		}
	}
}
