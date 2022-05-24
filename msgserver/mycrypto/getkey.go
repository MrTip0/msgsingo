package mycrypto

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	"net"
)

func GetKey(conn net.Conn) ([]byte, error) {
	prime, err := rand.Prime(rand.Reader, 64)
	if err != nil {
		return nil, fmt.Errorf("%s", err.Error())
	}
	num, err := rand.Int(rand.Reader, big.NewInt(10))
	if err != nil {
		return nil, fmt.Errorf("%s", err.Error())
	}

	fmt.Fprintf(conn, "%s;%s\x00", prime.String(), num.String())

	secret, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		return nil, fmt.Errorf("%s", err.Error())
	}

	othermod, err := bufio.NewReader(conn).ReadString('\x00')
	if err != nil {
		return nil, fmt.Errorf("%s", err.Error())
	}

	myhalf := new(big.Int)
	myhalf.Exp(num, secret, prime)
	_, err = fmt.Fprintf(conn, "%s\x00", myhalf.String())
	if err != nil {
		return nil, fmt.Errorf("%s", err.Error())
	}

	othnum := new(big.Int)

	othnum.SetString(othermod[:len(othermod)-1], 10)

	password := new(big.Int)

	password.Exp(othnum, secret, prime)

	return rightSizePass(password), nil
}
