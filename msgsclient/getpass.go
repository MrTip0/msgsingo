package main

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	"net"
	"os"
	"strings"
)

func getPass(conn net.Conn) []byte {
	password := new(big.Int)

inizializeconn:
	{
		prnum, err := bufio.NewReader(conn).ReadString('\x00')
		if err != nil {
			conn.Close()
			fmt.Printf("An error occurred: %s\n", err.Error())
			os.Exit(1)
		}

		if len(prnum) < 2 {
			goto inizializeconn
		}

		prnum = prnum[:len(prnum)-1]

		prnumarr := strings.Split(prnum, ";")
		prime, _ := new(big.Int).SetString(prnumarr[0], 10)
		number, _ := new(big.Int).SetString(prnumarr[1], 10)

		secret, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
		if err != nil {
			conn.Close()
			fmt.Printf("An error occurred: %s\n", err.Error())
			os.Exit(1)
		}

		myhalf := new(big.Int)
		myhalf.Exp(number, secret, prime)
		_, err = fmt.Fprintf(conn, "%s\x00", myhalf.String())
		if err != nil {
			conn.Close()
			fmt.Printf("An error occurred: %s\n", err.Error())
			os.Exit(1)
		}

		othermod, err := bufio.NewReader(conn).ReadString('\x00')
		if err != nil {
			conn.Close()
			fmt.Printf("An error occurred: %s\n", err.Error())
			os.Exit(1)
		}
		othnum := new(big.Int)

		othnum.SetString(othermod[:len(othermod)-1], 10)

		password.Exp(othnum, secret, prime)
	}

	return rightSizePass(password)
}
