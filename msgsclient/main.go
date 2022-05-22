package main

import (
	"bufio"
	"fmt"
	"math"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"golang.org/x/term"
)

type changetype uint8

const (
	Typed    changetype = 0
	Received            = 1
	Sent                = 2
)

type change struct {
	ctype changetype
	msg   string
}

func main() {
	var serverhost, username string
	{
		var err error
		reader := bufio.NewReader(os.Stdin)

		fmt.Printf("Insert the username that you want to use:\n> ")
		if username, err = reader.ReadString('\n'); err != nil {
			fmt.Printf("Error while reading input: %s", err.Error())
			os.Exit(1)
		}
		username = username[:len(username)-1]

		fmt.Printf("Insert the server ip:\n> ")
		if serverhost, err = reader.ReadString('\n'); err != nil {
			fmt.Printf("Error while reading input: %s", err.Error())
			os.Exit(1)
		}
		serverhost = serverhost[:len(serverhost)-1]
	}

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

	password := getPass(conn)

	{
		data, err := gcmEncrypter(password, username)
		if err != nil {
			fmt.Printf("Error while encripting username: %s\n", err.Error())
		}
		fmt.Fprintf(conn, "%s\x00", data)
	}

	end := make(chan bool, 1)

	sendchanges := make(chan change, 10)

	go sender(conn, password, sendchanges)
	go receiver(conn, end, password, sendchanges)
	go ui(sendchanges)

	<-end
}

func sender(conn net.Conn, pass []byte, updates chan change) {
	var msg string
	var err error

	reader := bufio.NewReader(os.Stdin)

	for {
		msg, err = readCharForChar(updates, reader)
		if err != nil {
			fmt.Printf("Error while reading input: %s\n", err.Error())
			continue
		}

		updates <- change{Sent, msg}

		data, err := gcmEncrypter(pass, msg)
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

func readCharForChar(updates chan change, reader *bufio.Reader) (string, error) {
	res := ""
	var tmp rune
	var err error
	cont := true
	for cont {
		if _, err = fmt.Scanf("%c", &tmp); err != nil {
			return "", err
		}
		if tmp != '\n' {
			strtmp := string(tmp)
			res += strtmp
			updates <- change{Typed, strtmp}
		} else {
			cont = false
		}
	}
	return res, nil
}

func receiver(conn net.Conn, end chan bool, pass []byte, updates chan change) {
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
			updates <- change{Received, data}
		}
	}
}

func drawUi(receivedsent []string, actmsg string, writer *bufio.Writer) {
	_, heigth, _ := term.GetSize(int(os.Stdin.Fd()))
	nscreen := make([]string, heigth)
	var usable []string

	usable = receivedsent[int(math.Max(float64(len(receivedsent)-heigth), 0)):]

	for i := 0; i < heigth-2 && i < len(usable); i++ {
		nscreen[i] = usable[i]
	}

	nscreen[heigth-1] = "> " + actmsg

	clearScreen()
	writer.Write([]byte(strings.Join(nscreen, "\n")))
	writer.Flush()
}

func ui(changes chan change) {
	receivedsent := make([]string, 0)
	actmsg := ""
	writer := bufio.NewWriter(os.Stdout)

	for {
		drawUi(receivedsent, actmsg, writer)
		ch := <-changes
		switch ch.ctype {
		case Sent:
			actmsg = ""
			receivedsent = append(receivedsent, "\t"+ch.msg)
		case Received:
			receivedsent = append(receivedsent, ch.msg)
		case Typed:
			actmsg += ch.msg
		}
	}
}

func clearScreen() {
	var clearcmd string = ""

	switch runtime.GOOS {
	case "windows":
		clearcmd = "cls"
	case "linux":
		clearcmd = "clear"
	case "darwin":
		clearcmd = "clear"
	}

	cmd := exec.Command(clearcmd)
	cmd.Stdout = os.Stdout
	cmd.Run()
}
