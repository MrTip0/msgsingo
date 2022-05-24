package mycrypto

import (
	"bufio"
	"fmt"
)

func ReadEncrypted(reader *bufio.Reader, key []byte) (res string, err error) {
begin:
	data, err := reader.ReadString('\x00')
	if err != nil {
		return
	}

	if len(data) > 0 {
		// Remove x00
		data = data[:len(data)-1]
		data, err = GcmDecrypter(key, data)
		if err != nil {
			return "", fmt.Errorf("error: %s\n", err.Error())
		}
	} else {
		goto begin
	}
	return data, nil
}
