package mycrypto

import (
	"fmt"
	"net"
)

func SendEncrypted(conn net.Conn, data string, key []byte) error {
	payload, err := GcmEncrypter(key, data)
	if err != nil {
		return fmt.Errorf("error: %s", err.Error())
	}

	_, err = fmt.Fprintf(conn, "%s\x00", payload)
	if err != nil {
		return fmt.Errorf("error: %s", err.Error())
	}

	return nil
}
